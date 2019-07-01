package statemachine

// EventCallbackFunc is a func with dynamic args. Any callback func of
// this type may accept a Transition object as input. Return values will be
// ignored.
//
// For AfterFailure callback, it must accept an `error` type arg:
//
// 	func(err error)
// 	func(transition eventmachine.Transition, err error)
type EventCallbackFunc interface{}

// EventCallbackBuilder provides the ability to define the `on` event(s)
// of the event callback matcher.
type EventCallbackBuilder interface {
	On(events ...string) EventCallbackOnBuilder
	OnAnyEvent() EventCallbackOnBuilder
	OnAnyEventExcept(events ...string) EventCallbackOnBuilder
}

// EventCallbackOnBuilder inherits from EventCallbackBuilder
// (or EventCallbackExceptFromBuilder) and provides the ability to define
// the transition callback func.
type EventCallbackOnBuilder interface {
	Do(callbackFunc EventCallbackFunc) EventCallbackOnBuilder
}

// newEventCallbackBuilder returns a zero-valued instance of
// eventCallbackBuilder, which implements
// EventCallbackBuilder.
func newEventCallbackBuilder(eventCallbackDef *EventCallbackDef) EventCallbackBuilder {
	return &eventCallbackBuilder{
		eventCallbackDef: eventCallbackDef,
	}
}

// eventCallbackBuilder implements EventCallbackBuilder
type eventCallbackBuilder struct {
	eventCallbackDef *EventCallbackDef
}

var _ EventCallbackBuilder = (*eventCallbackBuilder)(nil)

func (builder *eventCallbackBuilder) On(events ...string) EventCallbackOnBuilder {
	builder.eventCallbackDef.SetOn(events...)
	return newEventCallbackOnBuilder(builder.eventCallbackDef)
}

func (builder *eventCallbackBuilder) OnAnyEvent() EventCallbackOnBuilder {
	builder.eventCallbackDef.SetOnAnyEventExcept()
	return newEventCallbackOnBuilder(builder.eventCallbackDef)
}

func (builder *eventCallbackBuilder) OnAnyEventExcept(events ...string) EventCallbackOnBuilder {
	builder.eventCallbackDef.SetOnAnyEventExcept(events...)
	return newEventCallbackOnBuilder(builder.eventCallbackDef)
}

// newEventCallbackOnBuilder returns a zero-valued instance of
// eventCallbackOnBuilder, which implements
// EventCallbackOnBuilder.
func newEventCallbackOnBuilder(eventCallbackDef *EventCallbackDef) EventCallbackOnBuilder {
	return &eventCallbackOnBuilder{
		eventCallbackDef: eventCallbackDef,
	}
}

// eventCallbackOnBuilder implements EventCallbackOnBuilder
type eventCallbackOnBuilder struct {
	eventCallbackDef *EventCallbackDef
}

var _ EventCallbackOnBuilder = (*eventCallbackOnBuilder)(nil)

func (builder *eventCallbackOnBuilder) Do(callbackFunc EventCallbackFunc) EventCallbackOnBuilder {
	builder.eventCallbackDef.AddCallbackFunc(callbackFunc)
	return builder
}
