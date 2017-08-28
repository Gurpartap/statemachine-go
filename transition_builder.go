package statemachine

// TransitionGuardFunc may accept Transition object as input, and it must
// return a bool type.
//
// Valid TransitionGuardFunc signatures:
//
//	func() bool
//	func(transition statemachine.Transition) bool
type TransitionGuardFunc interface{}

// TransitionBuilder provides the ability to define the `from` state(s) of
// the transition matcher.
type TransitionBuilder interface {
	From(states ...string) TransitionFromBuilder
	FromAny() TransitionFromBuilder
	FromAnyExcept(states ...string) TransitionFromBuilder
}

// TransitionFromBuilder inherits `from` states from TransitionBuilder
// and provides the ability to define the `to` state as well as the
// `except from` states of the transition matcher.
type TransitionFromBuilder interface {
	ExceptFrom(states ...string) TransitionExceptFromBuilder
	To(state string) TransitionToBuilder
}

// TransitionExceptFromBuilder inherits from TransitionFromBuilder and
// provides the ability to define the `to` state of the transition matcher.
type TransitionExceptFromBuilder interface {
	To(state string) TransitionToBuilder
}

// TransitionToBuilder inherits from TransitionFromBuilder (or
// TransitionExceptFromBuilder) and provides the ability to define the guard
// condition funcs for the transition.
type TransitionToBuilder interface {
	If(guardFn TransitionGuardFunc)
	Unless(guardFn TransitionGuardFunc)
}

// newTransitionBuilder returns a zero-valued instance of
// TransitionBuilder, which implements TransitionBuilder.
func newTransitionBuilder() TransitionBuilder {
	return &transitionBuilder{}
}

// transitionBuilder implements TransitionBuilder
type transitionBuilder struct{}

func (builder *transitionBuilder) From(states ...string) TransitionFromBuilder {
	return newTransitionFromBuilder()
}

func (builder *transitionBuilder) FromAny() TransitionFromBuilder {
	return newTransitionFromBuilder()
}

func (builder *transitionBuilder) FromAnyExcept(states ...string) TransitionFromBuilder {
	return newTransitionFromBuilder()
}

// newTransitionFromBuilder returns a zero-valued instance of
// TransitionFromBuilder, which implements TransitionFromBuilder.
func newTransitionFromBuilder() TransitionFromBuilder {
	return &transitionFromBuilder{}
}

// transitionFromBuilder implements TransitionFromBuilder
type transitionFromBuilder struct{}

func (builder *transitionFromBuilder) ExceptFrom(states ...string) TransitionExceptFromBuilder {
	return newTransitionExceptFromBuilder()
}

func (builder *transitionFromBuilder) To(state string) TransitionToBuilder {
	return newTransitionToBuilder()
}

// newTransitionExceptFromBuilder returns a zero-valued instance of
// TransitionExceptFromBuilder, which implements TransitionExceptFromBuilder.
func newTransitionExceptFromBuilder() TransitionExceptFromBuilder {
	return &transitionExceptFromBuilder{}
}

// transitionExceptFromBuilder implements TransitionExceptFromBuilder
type transitionExceptFromBuilder struct{}

func (builder *transitionExceptFromBuilder) To(state string) TransitionToBuilder {
	return newTransitionToBuilder()
}

// newTransitionToBuilder returns a zero-valued instance of
// TransitionToBuilder, which implements TransitionToBuilder.
func newTransitionToBuilder() TransitionToBuilder {
	return &transitionToBuilder{}
}

// transitionToBuilder implements TransitionToBuilder
type transitionToBuilder struct{}

func (builder *transitionToBuilder) If(guardFn TransitionGuardFunc) {

}

func (builder *transitionToBuilder) Unless(guardFn TransitionGuardFunc) {

}
