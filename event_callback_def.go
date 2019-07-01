package statemachine

import (
	"fmt"
	"reflect"
)

type definedEvents struct {
	Events []string `json:",omitempty"`
	Except []string `json:",omitempty"`
}

type EventCallbackDef struct {
	On definedEvents
	Do []EventCallbackFunc `json:"-"`

	validateFor string
}

func (s *EventCallbackDef) MatchesEvent(event string) bool {
	for _, e := range s.On.Except {
		if e == event {
			return false
		}
	}

	if len(s.On.Events) == 0 {
		return true
	}

	for _, e := range s.On.Events {
		if e == event {
			return true
		}
	}

	return false
}

func (s *EventCallbackDef) SetOn(events ...string) {
	for _, event := range events {
		s.On.Events = append(s.On.Events, event)
	}
}

func (s *EventCallbackDef) SetOnAnyEventExcept(exceptEvents ...string) {
	for _, exceptEvent := range exceptEvents {
		s.On.Except = append(s.On.Except, exceptEvent)
	}
}

func (s *EventCallbackDef) AddCallbackFunc(callbackFunc EventCallbackFunc) {
	s.assertCallbackKind(callbackFunc)
	s.Do = append(s.Do, callbackFunc)
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
