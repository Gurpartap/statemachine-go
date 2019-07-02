package statemachine

import (
	"reflect"

	"github.com/Gurpartap/statemachine-go/internal/dynafunc"
)

type TransitionGuardDef struct {
	RegisteredFunc string          `json:",omitempty"`
	Guard          TransitionGuard `json:"-"`
}

type TransitionDef struct {
	From         []string `json:",omitempty"`
	ExceptFrom   []string `json:",omitempty"`
	To           string
	IfGuards     []*TransitionGuardDef `json:",omitempty"`
	UnlessGuards []*TransitionGuardDef `json:",omitempty"`
}

func execGuard(guard TransitionGuard, args map[reflect.Type]interface{}) bool {
	switch reflect.TypeOf(guard).Kind() {
	case reflect.Func:
		fn := dynafunc.NewDynamicFunc(guard, args)
		if err := fn.Call(); err != nil {
			panic(err)
		}
		// guard func must return a bool
		return fn.Out[0].Bool()
	case reflect.Ptr:
		if reflect.ValueOf(guard).Elem().Kind() == reflect.Bool {
			return reflect.ValueOf(guard).Elem().Bool()
		}
		fallthrough
	default:
		panic("guard must either be a compatible func or pointer to a bool variable")
	}
}

func (s *TransitionDef) IsAllowed(fromState string, machine Machine) bool {
	if len(s.IfGuards) != 0 || len(s.UnlessGuards) != 0 {
		args := make(map[reflect.Type]interface{})
		args[reflect.TypeOf(new(Machine))] = machine
		args[reflect.TypeOf(new(Transition))] = newTransitionImpl(
			fromState,
			s.To,
		)

		for _, guard := range s.IfGuards {
			// if !ok { dont allow }
			if ok := execGuard(guard.Guard, args); !ok {
				// fmt.Printf("❌1 from: %s to: %s\n", fromState, s.To.State())
				return false
			}
		}

		for _, guard := range s.UnlessGuards {
			// if ok { dont allow }
			if ok := execGuard(guard.Guard, args); ok {
				// fmt.Printf("❌2 from: %s to: %s\n", fromState, s.To.State())
				return false
			}
		}
	}

	// fmt.Printf("✅  transitioning from %s to %s\n", fromState, s.To.State())

	return true
}

func (s *TransitionDef) Matches(matchFrom string) bool {
	for _, exceptState := range s.ExceptFrom {
		if matchFrom == exceptState {
			return false
		}
	}

	// match any
	if len(s.From) == 0 {
		return true
	}

	for _, state := range s.From {
		if matchFrom == state {
			return true
		}
	}

	return false
}

func (s *TransitionDef) SetFrom(states ...string) {
	for _, state := range states {
		s.From = append(s.From, state)
	}
}

func (s *TransitionDef) SetFromAnyExcept(exceptStates ...string) {
	for _, exceptState := range exceptStates {
		s.ExceptFrom = append(s.ExceptFrom, exceptState)
	}
}

func (s *TransitionDef) SetTo(state string) {
	s.To = state
}

func (s *TransitionDef) AddIfGuard(guard TransitionGuard) {
	assertGuardKind(guard)
	s.IfGuards = append(s.IfGuards, &TransitionGuardDef{Guard: guard})
}

func (s *TransitionDef) AddUnlessGuard(guard TransitionGuard) {
	assertGuardKind(guard)
	s.UnlessGuards = append(s.UnlessGuards, &TransitionGuardDef{Guard: guard})
}

func assertGuardKind(guard TransitionGuard) {
	t := reflect.TypeOf(guard)
	switch t.Kind() {
	case reflect.Func:
		if t.NumIn() > 1 {
			panic("too many args in guard func")
		}
		if t.NumIn() == 1 && t.In(0).Implements(reflect.TypeOf(new(Transition)).Elem()) == false {
			panic("guard func arg must be statemachine.Transition")
		}
		if t.NumOut() != 1 {
			panic("guard func must return a single type")
		}
		if t.Out(0).Kind() != reflect.Bool {
			panic("guard func must return a bool type")
		}
		return
	case reflect.Ptr:
		if reflect.ValueOf(guard).Elem().Kind() == reflect.Bool {
			return
		}
	}
	panic("guard must either be a compatible func or pointer to a bool variable")
}
