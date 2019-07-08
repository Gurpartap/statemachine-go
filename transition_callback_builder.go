package statemachine

// TransitionCallbackFunc is a func with dynamic args. Any callback func of
// this type may accept a Machine and/or Transition object as inputs. Return
// values will be ignored.
//
// For BeforeTransition and AfterTransition:
//
// 	func()
// 	func(machine statemachine.Machine)
// 	func(transition statemachine.Transition)
// 	func(machine statemachine.Machine, transition statemachine.Transition)
//
// For AroundTransition callback, it must accept a `func()` type arg. The
// callback must call `next()` to continue the transition.
//
// 	func(next func())
// 	func(machine statemachine.Machine, next func())
// 	func(transition statemachine.Transition, next func())
// 	func(machine statemachine.Machine, transition statemachine.Transition, next func())
//
// For AfterFailure callback, it must accept an `error` type arg:
//
// 	func(err error)
// 	func(machine statemachine.Machine, err error)
// 	func(transition statemachine.Transition, err error)
// 	func(machine statemachine.Machine, transition statemachine.Transition, err error)
//
// TODO: perhaps support a Service interface struct, with methods to listen
//  for state changes. a service may be useful to implement [interruptible]
//  long-running callbacks. example: a download in a download manager.
type TransitionCallbackFunc interface{}

// TransitionCallbackBuilder provides the ability to define the `from`
// state(s) of the transition callback matcher.
type TransitionCallbackBuilder interface {
	From(states ...string) TransitionCallbackFromBuilder
	FromAny() TransitionCallbackFromBuilder
	FromAnyExcept(states ...string) TransitionCallbackFromBuilder
	To(states ...string) TransitionCallbackToBuilder
	ToAnyExcept(states ...string) TransitionCallbackToBuilder
	Any() TransitionCallbackToBuilder
}

// TransitionCallbackFromBuilder inherits `from` states from
// TransitionCallbackBuilder and provides the ability to define the
// `except from` and `to` states of the transition callback matcher.
type TransitionCallbackFromBuilder interface {
	ExceptFrom(states ...string) TransitionCallbackExceptFromBuilder
	To(states ...string) TransitionCallbackToBuilder
	ToSame() TransitionCallbackToBuilder
	ToAny() TransitionCallbackToBuilder
	ToAnyExcept(states ...string) TransitionCallbackToBuilder
	ToAnyExceptSame() TransitionCallbackToBuilder
}

// TransitionCallbackExceptFromBuilder inherits `from` states from
// TransitionCallbackBuilder and provides the ability to define the `to`
// states of the transition callback matcher.
type TransitionCallbackExceptFromBuilder interface {
	To(states ...string) TransitionCallbackToBuilder
	ToSame() TransitionCallbackToBuilder
	ToAny() TransitionCallbackToBuilder
	ToAnyExcept(states ...string) TransitionCallbackToBuilder
	ToAnyExceptSame() TransitionCallbackToBuilder
}

// TransitionCallbackToBuilder inherits from TransitionCallbackBuilder
// (or TransitionCallbackExceptFromBuilder) and provides the ability to define
// the transition callback func.
type TransitionCallbackToBuilder interface {
	ExitToState(supermachineState string)
	Do(callbackFuncs ...TransitionCallbackFunc) TransitionCallbackDoBuilder
}

type TransitionCallbackDoBuilder interface {
	Label(label string) TransitionCallbackToBuilder
}

// newTransitionCallbackBuilder returns a zero-valued instance of
// transitionCallbackBuilder, which implements
// TransitionCallbackBuilder.
func newTransitionCallbackBuilder(transitionCallbackDef *TransitionCallbackDef) TransitionCallbackBuilder {
	return &transitionCallbackBuilder{
		transitionCallbackDef: transitionCallbackDef,
	}
}

// transitionCallbackBuilder implements TransitionCallbackBuilder
type transitionCallbackBuilder struct {
	transitionCallbackDef *TransitionCallbackDef
}

var _ TransitionCallbackBuilder = (*transitionCallbackBuilder)(nil)

func (builder *transitionCallbackBuilder) From(states ...string) TransitionCallbackFromBuilder {
	builder.transitionCallbackDef.SetFrom(states...)
	return newTransitionCallbackFromBuilder(builder.transitionCallbackDef)
}

func (builder *transitionCallbackBuilder) FromAny() TransitionCallbackFromBuilder {
	builder.transitionCallbackDef.SetFromAnyExcept()
	return newTransitionCallbackFromBuilder(builder.transitionCallbackDef)
}

func (builder *transitionCallbackBuilder) FromAnyExcept(states ...string) TransitionCallbackFromBuilder {
	builder.transitionCallbackDef.SetFromAnyExcept(states...)
	return newTransitionCallbackFromBuilder(builder.transitionCallbackDef)
}

func (builder *transitionCallbackBuilder) To(states ...string) TransitionCallbackToBuilder {
	builder.transitionCallbackDef.SetFromAnyExcept()
	builder.transitionCallbackDef.SetTo(states...)
	return newTransitionCallbackToBuilder(builder.transitionCallbackDef)
}

func (builder *transitionCallbackBuilder) ToAnyExcept(states ...string) TransitionCallbackToBuilder {
	builder.transitionCallbackDef.SetFromAnyExcept()
	builder.transitionCallbackDef.SetToAnyExcept(states...)
	return newTransitionCallbackToBuilder(builder.transitionCallbackDef)
}

func (builder *transitionCallbackBuilder) Any() TransitionCallbackToBuilder {
	builder.transitionCallbackDef.SetFromAnyExcept()
	builder.transitionCallbackDef.SetToAnyExcept()
	return newTransitionCallbackToBuilder(builder.transitionCallbackDef)
}

// newTransitionCallbackFromBuilder returns a zero-valued instance of
// transitionCallbackFromBuilder, which implements
// TransitionCallbackFromBuilder.
func newTransitionCallbackFromBuilder(transitionCallbackDef *TransitionCallbackDef) TransitionCallbackFromBuilder {
	return &transitionCallbackFromBuilder{
		transitionCallbackDef: transitionCallbackDef,
	}
}

// transitionCallbackFromBuilder implements TransitionCallbackFromBuilder
type transitionCallbackFromBuilder struct {
	transitionCallbackDef *TransitionCallbackDef
}

var _ TransitionCallbackFromBuilder = (*transitionCallbackFromBuilder)(nil)

func (builder *transitionCallbackFromBuilder) ExceptFrom(states ...string) TransitionCallbackExceptFromBuilder {
	builder.transitionCallbackDef.SetFromAnyExcept(states...)
	return newTransitionCallbackExceptFromBuilder(builder.transitionCallbackDef)
}

func (builder *transitionCallbackFromBuilder) To(states ...string) TransitionCallbackToBuilder {
	builder.transitionCallbackDef.SetTo(states...)
	return newTransitionCallbackToBuilder(builder.transitionCallbackDef)
}

func (builder *transitionCallbackFromBuilder) ToSame() TransitionCallbackToBuilder {
	builder.transitionCallbackDef.SetSame()
	return newTransitionCallbackToBuilder(builder.transitionCallbackDef)
}

func (builder *transitionCallbackFromBuilder) ToAny() TransitionCallbackToBuilder {
	builder.transitionCallbackDef.SetToAnyExcept()
	return newTransitionCallbackToBuilder(builder.transitionCallbackDef)
}

func (builder *transitionCallbackFromBuilder) ToAnyExcept(states ...string) TransitionCallbackToBuilder {
	builder.transitionCallbackDef.SetToAnyExcept(states...)
	return newTransitionCallbackToBuilder(builder.transitionCallbackDef)
}

func (builder *transitionCallbackFromBuilder) ToAnyExceptSame() TransitionCallbackToBuilder {
	// builder.transitionCallbackDef.SetToAnyExceptSame()
	return newTransitionCallbackToBuilder(builder.transitionCallbackDef)
}

// newTransitionCallbackExceptFromBuilder returns a zero-valued instance of
// transitionCallbackExceptFromBuilder, which implements
// TransitionCallbackExceptFromBuilder.
func newTransitionCallbackExceptFromBuilder(transitionCallbackDef *TransitionCallbackDef) TransitionCallbackExceptFromBuilder {
	return &transitionCallbackExceptFromBuilder{
		transitionCallbackDef: transitionCallbackDef,
	}
}

// transitionCallbackExceptFromBuilder implements
// TransitionCallbackExceptFromBuilder
type transitionCallbackExceptFromBuilder struct {
	transitionCallbackDef *TransitionCallbackDef
}

var _ TransitionCallbackExceptFromBuilder = (*transitionCallbackExceptFromBuilder)(nil)

func (builder *transitionCallbackExceptFromBuilder) To(states ...string) TransitionCallbackToBuilder {
	builder.transitionCallbackDef.SetTo(states...)
	return newTransitionCallbackToBuilder(builder.transitionCallbackDef)
}

func (builder *transitionCallbackExceptFromBuilder) ToSame() TransitionCallbackToBuilder {
	// builder.transitionCallbackDef.SetSame()
	return newTransitionCallbackToBuilder(builder.transitionCallbackDef)
}

func (builder *transitionCallbackExceptFromBuilder) ToAny() TransitionCallbackToBuilder {
	builder.transitionCallbackDef.SetToAnyExcept()
	return newTransitionCallbackToBuilder(builder.transitionCallbackDef)
}

func (builder *transitionCallbackExceptFromBuilder) ToAnyExcept(states ...string) TransitionCallbackToBuilder {
	builder.transitionCallbackDef.SetToAnyExcept(states...)
	return newTransitionCallbackToBuilder(builder.transitionCallbackDef)
}

func (builder *transitionCallbackExceptFromBuilder) ToAnyExceptSame() TransitionCallbackToBuilder {
	// builder.transitionCallbackDef.SetToAnyExceptSame()
	return newTransitionCallbackToBuilder(builder.transitionCallbackDef)
}

// newTransitionCallbackToBuilder returns a zero-valued instance of
// transitionCallbackToBuilder, which implements
// TransitionCallbackToBuilder.
func newTransitionCallbackToBuilder(transitionCallbackDef *TransitionCallbackDef) TransitionCallbackToBuilder {
	return &transitionCallbackToBuilder{
		transitionCallbackDef: transitionCallbackDef,
	}
}

// transitionCallbackToBuilder implements TransitionCallbackToBuilder
type transitionCallbackToBuilder struct {
	transitionCallbackDef *TransitionCallbackDef
}

var _ TransitionCallbackToBuilder = (*transitionCallbackToBuilder)(nil)

func (builder *transitionCallbackToBuilder) ExitToState(supermachineState string) {
	builder.transitionCallbackDef.SetExitToState(supermachineState)
}

func (builder *transitionCallbackToBuilder) Do(callbackFuncs ...TransitionCallbackFunc) TransitionCallbackDoBuilder {
	builder.transitionCallbackDef.AddCallbackFunc(callbackFuncs...)
	return newTransitionCallbackDoBuilder(builder.transitionCallbackDef)
}

// newTransitionCallbackDoBuilder returns a zero-valued instance of
// transitionCallbackToBuilder, which implements
// TransitionCallbackDoBuilder.
func newTransitionCallbackDoBuilder(transitionCallbackDef *TransitionCallbackDef) TransitionCallbackDoBuilder {
	return &transitionCallbackDoBuilder{
		transitionCallbackDef: transitionCallbackDef,
	}
}

// transitionCallbackDoBuilder implements TransitionCallbackDoBuilder
type transitionCallbackDoBuilder struct {
	transitionCallbackDef *TransitionCallbackDef
}

var _ TransitionCallbackDoBuilder = (*transitionCallbackDoBuilder)(nil)

func (builder *transitionCallbackDoBuilder) Label(label string) TransitionCallbackToBuilder {
	builder.transitionCallbackDef.Do[len(builder.transitionCallbackDef.Do)-1].Label = label
	return newTransitionCallbackToBuilder(builder.transitionCallbackDef)
}
