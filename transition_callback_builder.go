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
type TransitionCallbackFunc interface{}

// TransitionCallbackBuilder provides the ability to define the `from`
// state(s) of the transition callback matcher.
type TransitionCallbackBuilder interface {
	From(states ...string) TransitionCallbackFromBuilder
	FromAny() TransitionCallbackFromBuilder
	FromAnyExcept(states ...string) TransitionCallbackFromBuilder
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
}

// TransitionCallbackExceptFromBuilder inherits `from` states from
// TransitionCallbackBuilder and provides the ability to define the `to`
// states of the transition callback matcher.
type TransitionCallbackExceptFromBuilder interface {
	To(states ...string) TransitionCallbackToBuilder
	ToSame() TransitionCallbackToBuilder
	ToAny() TransitionCallbackToBuilder
	ToAnyExcept(states ...string) TransitionCallbackToBuilder
}

// TransitionCallbackToBuilder inherits from TransitionCallbackBuilder
// (or TransitionCallbackExceptFromBuilder) and provides the ability to define
// the transition callback func.
type TransitionCallbackToBuilder interface {
	Do(callbackFunc TransitionCallbackFunc) TransitionCallbackToBuilder
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

func (builder *transitionCallbackToBuilder) Do(callbackFunc TransitionCallbackFunc) TransitionCallbackToBuilder {
	builder.transitionCallbackDef.AddCallbackFunc(callbackFunc)
	return builder
}
