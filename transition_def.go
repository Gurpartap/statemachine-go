package statemachine

import (
	"reflect"

	"github.com/Gurpartap/statemachine-go/internal/dynafunc"
)

type TransitionGuardDef struct {
	Label          string          `json:",omitempty" hcl:"label" hcle:"omitempty"`
	RegisteredFunc string          `json:",omitempty" hcl:"registered_func" hcle:"omitempty"`
	Guard          TransitionGuard `json:"-" hcle:"omit"`
}

type TransitionDef struct {
	From         []string              `json:",omitempty" hcl:"from" hcle:"omitempty"`
	ExceptFrom   []string              `json:",omitempty" hcl:"except_from" hcle:"omitempty"`
	To           string                `hcl:"to"`
	IfGuards     []*TransitionGuardDef `json:",omitempty" hcl:"if_guard" hcle:"omitempty"`
	UnlessGuards []*TransitionGuardDef `json:",omitempty" hcl:"unless_guard" hcle:"omitempty"`
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

func (def *TransitionDef) IsAllowed(fromState string, machine Machine) bool {
	if len(def.IfGuards) != 0 || len(def.UnlessGuards) != 0 {
		args := make(map[reflect.Type]interface{})
		args[reflect.TypeOf(new(Machine))] = machine
		args[reflect.TypeOf(new(Transition))] = newTransitionImpl(
			fromState,
			def.To,
		)

		for _, guard := range def.IfGuards {
			// if !ok { dont allow }
			if ok := execGuard(guard.Guard, args); !ok {
				// fmt.Printf("❌1 from: %def to: %def\n", fromState, def.To.State())
				return false
			}
		}

		for _, guard := range def.UnlessGuards {
			// if ok { dont allow }
			if ok := execGuard(guard.Guard, args); ok {
				// fmt.Printf("❌2 from: %def to: %def\n", fromState, def.To.State())
				return false
			}
		}
	}

	// fmt.Printf("✅  transitioning from %def to %def\n", fromState, def.To.State())

	return true
}

func (def *TransitionDef) Matches(matchFrom string) bool {
	for _, exceptState := range def.ExceptFrom {
		if matchFrom == exceptState {
			return false
		}
	}

	// match any
	if len(def.From) == 0 {
		return true
	}

	for _, state := range def.From {
		if matchFrom == state {
			return true
		}
	}

	return false
}

func (def *TransitionDef) SetFrom(states ...string) {
	for _, state := range states {
		def.From = append(def.From, state)
	}
}

func (def *TransitionDef) SetFromAnyExcept(exceptStates ...string) {
	for _, exceptState := range exceptStates {
		def.ExceptFrom = append(def.ExceptFrom, exceptState)
	}
}

func (def *TransitionDef) SetTo(state string) {
	def.To = state
}

func (def *TransitionDef) AddIfGuard(guards ...TransitionGuard) {
	for _, guard := range guards {
		assertGuardKind(guard)
		def.IfGuards = append(def.IfGuards, &TransitionGuardDef{Guard: guard})
	}
}

func (def *TransitionDef) AddUnlessGuard(guards ...TransitionGuard) {
	for _, guard := range guards {
		assertGuardKind(guard)
		def.UnlessGuards = append(def.UnlessGuards, &TransitionGuardDef{Guard: guard})
	}
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
