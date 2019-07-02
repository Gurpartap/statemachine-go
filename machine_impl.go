package statemachine

import (
	"reflect"
	"sync"

	"github.com/Gurpartap/statemachine-go/internal/dynafunc"
)

type machineImpl struct {
	def *MachineDef

	previousState string
	currentState  string

	mutex sync.RWMutex
}

// NewMachine returns a zero-valued instance of machine, which implements
// Machine.
func NewMachine() Machine {
	return &machineImpl{
		def:         NewMachineDef(),
	}
}

// BuildNewMachine creates a zero-valued instance of machine, and builds it
// using the passed machineBuilderFn arg.
func BuildNewMachine(machineBuilderFn func(machineBuilder MachineBuilder)) Machine {
	machine := NewMachine()
	machine.Build(machineBuilderFn)
	return machine
}

func (m *machineImpl) Build(machineBuilderFn func(machineBuilder MachineBuilder)) {
	machineBuilder := NewMachineBuilder()
	machineBuilderFn(machineBuilder)
	machineBuilder.Build(m)
}

// SetMachineDef implements MachineBuildable.
func (m *machineImpl) SetMachineDef(def *MachineDef) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// b, _ := json.MarshalIndent(def, "", "  ")
	// fmt.Printf("machine def = %s\n", string(b))

	m.def = def
	m.setCurrentState(m.def.InitialState)
}

// GetState implements Machine.
func (m *machineImpl) GetState() string {
	return m.currentState
}

// SetCurrentState implements Machine.
func (m *machineImpl) SetCurrentState(state string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.setCurrentState(state)
}

// IsState implements Machine.
func (m *machineImpl) IsState(state string) bool {
	return m.GetState() == state
}

// Fire implements Machine.
func (m *machineImpl) Fire(event string) (err error) {
	m.mutex.Lock()
	defer func() {
		if err != nil {
			args := make(map[reflect.Type]interface{})
			args[reflect.TypeOf(new(Event))] = &eventImpl{name: event}
			args[reflect.TypeOf(new(error))] = err

			for _, callbackDef := range m.def.FailureCallbacks {
				if callbackDef.MatchesEvent(event) {
					for _, callback := range callbackDef.Do {
						m.exec(callback.Func, args)
					}
				}
			}
		}

		m.mutex.Unlock()
	}()

	if m.IsState("") {
		err = ErrNotInitialized
		return
	}

	// fmt.Printf("\n---\n🔁 %s\n", event)
	// defer func() { fmt.Printf("=> %s\n---\n", m.GetState()) }()

	// if transitionDef, ok := m.def.match(); ok {}

	fromState := m.GetState()

	eventDef, ok := m.def.Events[event]
	if !ok {
		err = ErrNoSuchEvent
		return
	}

	var transition Transition
	err = ErrNoMatchingTransition

	for _, transitionDef := range eventDef.Transitions {
		matches := transitionDef.Matches(fromState)
		if !matches {
			err = ErrNoMatchingTransition
			continue
		}
		if !transitionDef.IsAllowed(fromState) {
			err = ErrTransitionNotAllowed
			continue
		}

		transition = newTransitionImpl(fromState, transitionDef.To)
		err = nil

		break
	}
	if err != nil {
		return err
	}

	args := make(map[reflect.Type]interface{})
	args[reflect.TypeOf(new(Transition))] = transition

	for _, callbackDef := range m.def.BeforeCallbacks {
		if callbackDef.Matches(fromState, transition.To()) {
			for _, callback := range callbackDef.Do {
				m.exec(callback.Func, args)
			}
		}
	}

	var matchingCallbacks []*TransitionCallbackFuncDef
	for _, callbackDef := range m.def.AroundCallbacks {
		if callbackDef.Matches(fromState, transition.To()) {
			matchingCallbacks = append(matchingCallbacks, callbackDef.Do...)
		}
	}
	applyTransition := func() {
		m.setCurrentState(transition.To())
	}

	m.applyTransitionAroundCallbacks(matchingCallbacks, args, applyTransition)

	for _, callbackDef := range m.def.AfterCallbacks {
		if callbackDef.Matches(fromState, transition.To()) {
			for _, callback := range callbackDef.Do {
				m.exec(callback.Func, args)
			}
		}
	}

	return nil
}

func (m *machineImpl) setCurrentState(state string) {
	for _, s := range m.def.States {
		if s == state {
			m.previousState = m.currentState
			m.currentState = state
			break
		}
	}
}

// callback1(next: {
//   callback2(next: {
//     callback3(next: {
//       applyTransition()
//     })
//   })
// })
func (m *machineImpl) applyTransitionAroundCallbacks(callbacks []*TransitionCallbackFuncDef, args map[reflect.Type]interface{}, applyTransition func()) {
	if len(callbacks) == 0 {
		applyTransition()
		return
	}

	calledBackNext := false

	args[reflect.TypeOf(new(func()))] = func() {
		calledBackNext = true
		m.applyTransitionAroundCallbacks(callbacks[1:], args, applyTransition)
	}

	m.exec(callbacks[0].Func, args)
	if !calledBackNext && len(callbacks) != 1 {
		panic("non-last around callbacks must call next()")
	}

	return
}

func (m *machineImpl) exec(callback TransitionCallbackFunc, args map[reflect.Type]interface{}) {
	fn := dynafunc.NewDynamicFunc(callback, args)
	if err := fn.Call(); err != nil {
		panic(err)
	}
}
