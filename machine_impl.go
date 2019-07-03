package statemachine

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/Gurpartap/statemachine-go/internal/dynafunc"
)

type machineImpl struct {
	def *MachineDef

	previousState string
	currentState  string

	supermachine *machineImpl
	submachines  map[string][]*machineImpl

	mutex     sync.RWMutex
	hasExited bool

	ctxTimedEvents  context.Context
	stopTimedEvents context.CancelFunc
}

// NewMachine returns a zero-valued instance of machine, which implements
// Machine.
func NewMachine() Machine {
	ctxTimedEvents, stopTimedEvents := context.WithCancel(context.Background())
	return &machineImpl{
		def:             NewMachineDef(),
		submachines:     map[string][]*machineImpl{},
		ctxTimedEvents:  ctxTimedEvents,
		stopTimedEvents: stopTimedEvents,
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
	// // b, _ := hclencoder.Encode(def)
	// fmt.Printf("machine def = %s\n", string(b))

	m.def = def
	if err := m.setCurrentState(m.def.InitialState); err != nil {
		panic(err)
	}
	m.restartTimedEventsLoops()
}

func (m *machineImpl) restartTimedEventsLoops() {
	for event, eventDef := range m.def.Events {
		if eventDef.TimedEvery > 0 {
			go func(event string, timedEvery time.Duration) {
				// fmt.Printf("event=%s timed_every=%d\n", event, timedEvery)
				for {
					select {
					case <-time.After(timedEvery):
						// fmt.Printf("firing timed event '%s'\n", event)
						_ = m.Fire(event)
					case <-m.ctxTimedEvents.Done():
						// fmt.Printf("stopping timed event '%s'\n", event)
						return
					}
				}
			}(event, eventDef.TimedEvery)
		}
	}
}

// GetStateMap implements Machine.
func (m *machineImpl) GetStateMap() StateMap {
	substate := StateMap{}
	if submachines, ok := m.submachines[m.currentState]; ok {
		for _, submachine := range submachines {
			if _, ok := submachine.submachines[submachine.currentState]; ok {
				substate[submachine.def.ID] = submachine.GetStateMap()
			} else {
				substate[submachine.def.ID] = submachine.currentState
			}
		}
	}
	return StateMap{
		m.currentState: substate,
	}
}

// GetState implements Machine.
func (m *machineImpl) GetState() string {
	return m.currentState
}

// SetCurrentState implements Machine.
func (m *machineImpl) SetCurrentState(state interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.setCurrentState(state)
}

// IsState implements Machine.
func (m *machineImpl) IsState(state string) bool {
	return m.GetState() == state
}

// Send implements Machine.
func (m *machineImpl) Send(signal Message) error {
	switch signal.(type) {
	case TriggerEvent:
		return m.Fire(signal.(TriggerEvent).Event)
	case OverrideState:
		return m.SetCurrentState(signal.(OverrideState).State)
	}
	return errors.New("no such signal")
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
		if m.hasExited {
			// TODO: should we wait for `<-m.stoppedTimedEvents`?
			m.stopTimedEvents()
			*m = machineImpl{}
		}
	}()

	if m.IsState("") {
		err = errors.New("state machine not initialized")
		return
	}

	// fmt.Printf("\n---\n🔁 %s\n", event)
	// defer func() { fmt.Printf("=> %s\n---\n", m.GetState()) }()

	fromState := m.GetState()

	var transition Transition
	transition, err = m.findTransition(event, fromState)
	if err != nil {
		return
	}

	err = m.applyTransition(transition)
	return
}

func (m *machineImpl) findTransition(event string, fromState string) (transition Transition, err error) {
	eventDef, ok := m.def.Events[event]
	if !ok {
		err = errors.New("no such event")
		return
	}

	transition, err = m.matchTransition(eventDef.Transitions, fromState)
	if err == nil || eventDef.Choice == nil {
		return
	}

	transition, err = m.findChoiceTransition(event, eventDef, fromState)
	return
}

func (m *machineImpl) findChoiceTransition(event string, eventDef *EventDef, fromState string) (transition Transition, err error) {
	args := make(map[reflect.Type]interface{})
	args[reflect.TypeOf(new(Machine))] = m
	args[reflect.TypeOf(new(Event))] = &eventImpl{name: event}

	if eventDef.Choice.UnlessGuard != nil {
		if ok := execGuard(eventDef.Choice.UnlessGuard.Guard, args); ok {
			err = ErrTransitionNotAllowed
			return
		}
	}

	if execChoice(eventDef.Choice.Condition.Condition, args) {
		if eventDef.Choice.OnTrue.Choice != nil {
			transition, err = m.findChoiceTransition(event, eventDef.Choice.OnTrue, fromState)
			return
		}
		transition, err = m.matchTransition(eventDef.Choice.OnTrue.Transitions, fromState)
		return
	}

	if eventDef.Choice.OnFalse.Choice != nil {
		transition, err = m.findChoiceTransition(event, eventDef.Choice.OnFalse, fromState)
		return
	}
	transition, err = m.matchTransition(eventDef.Choice.OnFalse.Transitions, fromState)
	return
}

func (m *machineImpl) matchTransition(transitions []*TransitionDef, fromState string) (transition Transition, err error) {
	for _, transitionDef := range transitions {
		matches := transitionDef.Matches(fromState)
		if !matches {
			err = ErrNoMatchingTransition
			continue
		}
		if !transitionDef.IsAllowed(fromState, m) {
			err = ErrTransitionNotAllowed
			continue
		}

		transition = newTransitionImpl(fromState, transitionDef.To)
		err = nil

		return
	}

	err = ErrNoMatchingTransition
	return
}

func (m *machineImpl) Submachine(idPath ...string) (Machine, error) {
	for _, submachine := range m.submachines[m.currentState] {
		if submachine.def.ID == idPath[0] {
			if len(idPath) > 1 {
				return submachine.Submachine(idPath[1:]...)
			}
			return submachine, nil
		}
	}

	return nil, errors.New("submachine not active")
}

func (m *machineImpl) setCurrentStateMap(state StateMap) error {
	for rootState, subStates := range state {
		switch subStates.(type) {
		case nil:
			// fmt.Printf("setting state to '%s'\n", rootState)
			return m.setCurrentState(rootState)

		case StateMap:
			// fmt.Printf("setting state to '%s'\n", rootState)
			if err := m.setCurrentState(rootState); err != nil {
				return err
			}

			for id, state := range subStates.(StateMap) {
				for _, submachine := range m.submachines[rootState] {
					if submachine.def.ID == id {
						switch state.(type) {
						case StateMap:
							// fmt.Printf("nesting into submachine '%s'\n", id)
							if err := submachine.setCurrentStateMap(state.(StateMap)); err != nil {
								return err
							}
						case string:
							// fmt.Printf("setting submachine '%s' to '%s'\n", id, state)
							if err := submachine.SetCurrentState(state); err != nil {
								return err
							}
						default:
							return ErrStateTypeNotSupported
						}
					}
				}
			}

		default:
			if err := m.setCurrentState(rootState); err != nil {
				return err
			}

			return ErrStateTypeNotSupported
		}

		// there is only one kv in any given StateMap
		break
	}

	return nil
}

func (m *machineImpl) setCurrentState(state interface{}) error {
	if state, ok := state.(StateMap); ok {
		if err := m.setCurrentStateMap(state); err != nil {
			return err
		}
		return nil
	}

	if state, ok := state.(string); ok {
		for _, s := range m.def.States {
			if s == state {
				m.previousState = m.currentState
				m.currentState = state
				return nil
			}
		}

		for s, submachineDefs := range m.def.Submachines {
			if s == state {
				m.submachines[state] = []*machineImpl{}
				for _, submachineDef := range submachineDefs {
					submachine := &machineImpl{
						supermachine: m,
						submachines:  map[string][]*machineImpl{},
					}
					submachine.SetMachineDef(submachineDef)
					m.submachines[state] = append(m.submachines[state], submachine)
				}

				m.previousState = m.currentState
				m.currentState = state
				return nil
			}
		}
	}

	return ErrStateTypeNotSupported
}

func (m *machineImpl) applyTransition(transition Transition) error {
	fromState := m.GetState()

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
		if !callbackDef.Matches(fromState, transition.To()) {
			continue
		}
		for _, callback := range callbackDef.Do {
			m.exec(callback.Func, args)
		}
		if callbackDef.ExitToState != "" && m.supermachine != nil {
			if err := m.supermachine.applyTransition(
				newTransitionImpl(m.supermachine.currentState, callbackDef.ExitToState),
			); err != nil {
				return fmt.Errorf("could not exit submachine: %s", err)
			}
			m.hasExited = true
			return nil
		}
	}

	return nil
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
	args[reflect.TypeOf(new(Machine))] = m
	fn := dynafunc.NewDynamicFunc(callback, args)
	if err := fn.Call(); err != nil {
		panic(err)
	}
}
