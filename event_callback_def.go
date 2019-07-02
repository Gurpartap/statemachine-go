package statemachine

import (
	"fmt"
	"reflect"
)

type EventCallbackFuncDef struct {
	RegisteredFunc string            `json:",omitempty"`
	Func           EventCallbackFunc `json:"-"`
}

type EventCallbackDef struct {
	On       []string                `json:",omitempty"`
	ExceptOn []string                `json:",omitempty"`
	Do       []*EventCallbackFuncDef `json:",omitempty"`

	validateFor string `json:"-"`
}

func (s *EventCallbackDef) MatchesEvent(event string) bool {
	for _, e := range s.ExceptOn {
		if e == event {
			return false
		}
	}

	if len(s.On) == 0 {
		return true
	}

	for _, e := range s.On {
		if e == event {
			return true
		}
	}

	return false
}

func (s *EventCallbackDef) SetOn(events ...string) {
	for _, event := range events {
		s.On = append(s.On, event)
	}
}

func (s *EventCallbackDef) SetOnAnyEventExcept(exceptEvents ...string) {
	for _, exceptEvent := range exceptEvents {
		s.ExceptOn = append(s.ExceptOn, exceptEvent)
	}
}

func (s *EventCallbackDef) AddCallbackFunc(callbackFunc EventCallbackFunc) {
	s.assertCallbackKind(callbackFunc)
	s.Do = append(s.Do, &EventCallbackFuncDef{Func: callbackFunc})
}

func (s *EventCallbackDef) assertCallbackKind(callbackFunc EventCallbackFunc) {
	t := reflect.TypeOf(callbackFunc)
	switch t.Kind() {
	case reflect.Func:
		if t.NumOut() != 0 {
			panic("callback func must not return anything")
		}

		optionalArgs := make(map[reflect.Type]struct{})
		requiredArgs := make(map[reflect.Type]struct{})

		switch s.validateFor {
		case "AfterFailure":
			optionalArgs[reflect.TypeOf(new(Event))] = struct{}{}
			requiredArgs[reflect.TypeOf(new(error))] = struct{}{}
		}

		// ensure all args are of expected types, whether optional or required
		for i := 0; i < t.NumIn(); i++ {
			argType := t.In(i)
			if _, ok := optionalArgs[reflect.PtrTo(argType)]; ok {
				continue
			}
			if _, ok := requiredArgs[reflect.PtrTo(argType)]; ok {
				continue
			}
			panic(fmt.Sprintf("unexpected argument with type '%s' in %s callback", argType, s.validateFor))
		}

	outer:
		// ensure required args are present
		for requiredArg := range requiredArgs {
			for i := 0; i < t.NumIn(); i++ {
				if requiredArg.Elem() == t.In(i) {
					continue outer
				}
			}

			panic(fmt.Sprintf("missing required arg '%s' in %s callback", requiredArg.Elem().String(), s.validateFor))
		}
		return
	}
	panic("callback must be a compatible func")
}
