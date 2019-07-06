package statemachine

// TransitionGuard may accept Transition object as input, and it must
// return a bool type.
//
// Valid TransitionGuard types:
//
//  bool
// 	func() bool
// 	func(transition statemachine.Transition) bool
type TransitionGuard interface{}

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
	If(guards ...TransitionGuard) TransitionAndGuardBuilder
	Unless(guards ...TransitionGuard) TransitionAndGuardBuilder
}

// TransitionAndGuardBuilder inherits from TransitionToBuilder and provides
// the ability to define additional guard condition funcs for the transition.
type TransitionAndGuardBuilder interface {
	Label(label string) TransitionAndGuardBuilder
	AndIf(guards ...TransitionGuard) TransitionAndGuardBuilder
	AndUnless(guards ...TransitionGuard) TransitionAndGuardBuilder
}

// newTransitionBuilder returns a zero-valued instance of
// TransitionBuilder, which implements TransitionBuilder.
func newTransitionBuilder(transitionDef *TransitionDef) TransitionBuilder {
	return &transitionBuilder{
		transitionDef: transitionDef,
	}
}

// transitionBuilder implements TransitionBuilder
type transitionBuilder struct {
	transitionDef *TransitionDef
}

var _ TransitionBuilder = (*transitionBuilder)(nil)

func (builder *transitionBuilder) From(states ...string) TransitionFromBuilder {
	builder.transitionDef.SetFrom(states...)
	return newTransitionFromBuilder(builder.transitionDef)
}

func (builder *transitionBuilder) FromAny() TransitionFromBuilder {
	builder.transitionDef.SetFromAnyExcept()
	return newTransitionFromBuilder(builder.transitionDef)
}

func (builder *transitionBuilder) FromAnyExcept(states ...string) TransitionFromBuilder {
	builder.transitionDef.SetFromAnyExcept(states...)
	return newTransitionFromBuilder(builder.transitionDef)
}

// newTransitionFromBuilder returns a zero-valued instance of
// TransitionFromBuilder, which implements TransitionFromBuilder.
func newTransitionFromBuilder(transitionDef *TransitionDef) TransitionFromBuilder {
	return &transitionFromBuilder{
		transitionDef: transitionDef,
	}
}

// transitionFromBuilder implements TransitionFromBuilder
type transitionFromBuilder struct {
	transitionDef *TransitionDef
}

var _ TransitionFromBuilder = (*transitionFromBuilder)(nil)

func (builder *transitionFromBuilder) ExceptFrom(states ...string) TransitionExceptFromBuilder {
	builder.transitionDef.SetFromAnyExcept(states...)
	return newTransitionExceptFromBuilder(builder.transitionDef)
}

func (builder *transitionFromBuilder) To(state string) TransitionToBuilder {
	builder.transitionDef.SetTo(state)
	return newTransitionToBuilder(builder.transitionDef)
}

// newTransitionExceptFromBuilder returns a zero-valued instance of
// TransitionExceptFromBuilder, which implements TransitionExceptFromBuilder.
func newTransitionExceptFromBuilder(transitionDef *TransitionDef) TransitionExceptFromBuilder {
	return &transitionExceptFromBuilder{
		transitionDef: transitionDef,
	}
}

// transitionExceptFromBuilder implements TransitionExceptFromBuilder
type transitionExceptFromBuilder struct {
	transitionDef *TransitionDef
}

var _ TransitionExceptFromBuilder = (*transitionExceptFromBuilder)(nil)

func (builder *transitionExceptFromBuilder) To(state string) TransitionToBuilder {
	builder.transitionDef.SetTo(state)
	return newTransitionToBuilder(builder.transitionDef)
}

// newTransitionToBuilder returns a zero-valued instance of
// TransitionToBuilder, which implements TransitionToBuilder.
func newTransitionToBuilder(transitionDef *TransitionDef) TransitionToBuilder {
	return &transitionToBuilder{
		transitionDef: transitionDef,
	}
}

// transitionToBuilder implements TransitionToBuilder
type transitionToBuilder struct {
	transitionDef *TransitionDef
}

var _ TransitionToBuilder = (*transitionToBuilder)(nil)

func (builder *transitionToBuilder) If(guard ...TransitionGuard) TransitionAndGuardBuilder {
	builder.transitionDef.AddIfGuard(guard...)
	return newTransitionAndGuardBuilder(builder.transitionDef, "if")
}

func (builder *transitionToBuilder) Unless(guard ...TransitionGuard) TransitionAndGuardBuilder {
	builder.transitionDef.AddUnlessGuard(guard...)
	return newTransitionAndGuardBuilder(builder.transitionDef, "unless")
}

// newTransitionAndGuardBuilder returns a zero-valued instance of
// TransitionAndGuardBuilder, which implements TransitionAndGuardBuilder.
func newTransitionAndGuardBuilder(transitionDef *TransitionDef, lastGuardType string) TransitionAndGuardBuilder {
	return &transitionAndGuardBuilder{
		transitionDef: transitionDef,
		lastGuardType: lastGuardType,
	}
}

// transitionAndGuardBuilder implements TransitionAndGuardBuilder
type transitionAndGuardBuilder struct {
	transitionDef *TransitionDef
	lastGuardType string
}

var _ TransitionAndGuardBuilder = (*transitionAndGuardBuilder)(nil)

func (builder *transitionAndGuardBuilder) Label(label string) TransitionAndGuardBuilder {
	if builder.lastGuardType == "if" {
		builder.transitionDef.IfGuards[len(builder.transitionDef.IfGuards)-1].Label = label
	} else {
		builder.transitionDef.UnlessGuards[len(builder.transitionDef.UnlessGuards)-1].Label = label
	}
	return builder
}

func (builder *transitionAndGuardBuilder) AndIf(guard ...TransitionGuard) TransitionAndGuardBuilder {
	builder.transitionDef.AddIfGuard(guard...)
	builder.lastGuardType = "if"
	return builder
}

func (builder *transitionAndGuardBuilder) AndUnless(guard ...TransitionGuard) TransitionAndGuardBuilder {
	builder.transitionDef.AddUnlessGuard(guard...)
	builder.lastGuardType = "unless"
	return builder
}
