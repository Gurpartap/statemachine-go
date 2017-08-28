package statemachine

// TransitionCallbackFunc is a func with dynamic args. Any callback func of
// this type may accept a Transition object as input. Return values will be
// ignored.
//
// For BeforeTransition and AfterTransition:
//
//	func()
//	func(t statemachine.Transition)
//
// For AroundTransition callback, it must accept a `func` type arg:
//
//	func(execFn func())
//	func(t statemachine.Transition, execFn func())
//
// For AfterFailure callback, the func may optionally accept an `error` type
// arg:
//
//	func()
//	func(err error)
//	func(t statemachine.Transition, err error)
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
	Do(callbackFn TransitionCallbackFunc)
}

// newTransitionCallbackBuilder returns a zero-valued instance of
// transitionCallbackBuilder, which implements
// TransitionCallbackBuilder.
func newTransitionCallbackBuilder() TransitionCallbackBuilder {
	return &transitionCallbackBuilder{}
}

// transitionCallbackBuilder implements TransitionCallbackBuilder
type transitionCallbackBuilder struct{}

func (builder *transitionCallbackBuilder) From(states ...string) TransitionCallbackFromBuilder {
	return newTransitionCallbackFromBuilder()
}

func (builder *transitionCallbackBuilder) FromAny() TransitionCallbackFromBuilder {
	return newTransitionCallbackFromBuilder()
}

func (builder *transitionCallbackBuilder) FromAnyExcept(states ...string) TransitionCallbackFromBuilder {
	return newTransitionCallbackFromBuilder()
}

// newTransitionCallbackFromBuilder returns a zero-valued instance of
// transitionCallbackFromBuilder, which implements
// TransitionCallbackFromBuilder.
func newTransitionCallbackFromBuilder() TransitionCallbackFromBuilder {
	return &transitionCallbackFromBuilder{}
}

// transitionCallbackFromBuilder implements TransitionCallbackFromBuilder
type transitionCallbackFromBuilder struct{}

func (builder *transitionCallbackFromBuilder) ExceptFrom(states ...string) TransitionCallbackExceptFromBuilder {
	return newTransitionCallbackExceptFromBuilder()
}

func (builder *transitionCallbackFromBuilder) To(states ...string) TransitionCallbackToBuilder {
	return newTransitionCallbackToBuilder()
}

func (builder *transitionCallbackFromBuilder) ToSame() TransitionCallbackToBuilder {
	return newTransitionCallbackToBuilder()
}

func (builder *transitionCallbackFromBuilder) ToAny() TransitionCallbackToBuilder {
	return newTransitionCallbackToBuilder()
}

func (builder *transitionCallbackFromBuilder) ToAnyExcept(states ...string) TransitionCallbackToBuilder {
	return newTransitionCallbackToBuilder()
}

// newTransitionCallbackExceptFromBuilder returns a zero-valued instance of
// transitionCallbackExceptFromBuilder, which implements
// TransitionCallbackExceptFromBuilder.
func newTransitionCallbackExceptFromBuilder() TransitionCallbackExceptFromBuilder {
	return &transitionCallbackExceptFromBuilder{}
}

// transitionCallbackExceptFromBuilder implements
// TransitionCallbackExceptFromBuilder
type transitionCallbackExceptFromBuilder struct{}

func (builder *transitionCallbackExceptFromBuilder) To(states ...string) TransitionCallbackToBuilder {
	return newTransitionCallbackToBuilder()
}

func (builder *transitionCallbackExceptFromBuilder) ToSame() TransitionCallbackToBuilder {
	return newTransitionCallbackToBuilder()
}

func (builder *transitionCallbackExceptFromBuilder) ToAny() TransitionCallbackToBuilder {
	return newTransitionCallbackToBuilder()
}

func (builder *transitionCallbackExceptFromBuilder) ToAnyExcept(states ...string) TransitionCallbackToBuilder {
	return newTransitionCallbackToBuilder()
}

// newTransitionCallbackToBuilder returns a zero-valued instance of
// transitionCallbackToBuilder, which implements
// TransitionCallbackToBuilder.
func newTransitionCallbackToBuilder() TransitionCallbackToBuilder {
	return &transitionCallbackToBuilder{}
}

// transitionCallbackToBuilder implements TransitionCallbackToBuilder
type transitionCallbackToBuilder struct{}

func (builder *transitionCallbackToBuilder) Do(callbackFn TransitionCallbackFunc) {

}
