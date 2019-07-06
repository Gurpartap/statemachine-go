package statemachine

// ChoiceCondition may accept Transition object as input, and it must
// return a bool type.
//
// Valid ChoiceCondition types:
//
//  bool
// 	func() bool
// 	func(transition statemachine.Transition) bool
type ChoiceCondition interface{}

// ChoiceBuilder provides the ability to define the conditions and result
// handling of the choice definition.
type ChoiceBuilder interface {
	Label(label string) ChoiceBuilder
	Unless(guard TransitionGuard) ChoiceBuilder
	OnTrue(eventBuilderFn func(eventBuilder EventBuilder)) ChoiceTrueBuilder
	OnFalse(eventBuilderFn func(eventBuilder EventBuilder)) ChoiceFalseBuilder
}

// ChoiceTrueBuilder inherits from ChoiceBuilder
type ChoiceTrueBuilder interface {
	OnFalse(eventBuilderFn func(eventBuilder EventBuilder))
}

// ChoiceFalseBuilder inherits from ChoiceBuilder
type ChoiceFalseBuilder interface {
	OnTrue(eventBuilderFn func(eventBuilder EventBuilder))
}

// newChoiceBuilder returns a zero-valued instance of
// ChoiceBuilder, which implements ChoiceBuilder.
func newChoiceBuilder(choiceDef *ChoiceDef) ChoiceBuilder {
	return &chosenBuilder{
		choiceDef: choiceDef,
	}
}

// chosenBuilder implements ChoiceBuilder
type chosenBuilder struct {
	choiceDef *ChoiceDef
}

func (builder *chosenBuilder) Label(label string) ChoiceBuilder {
	builder.choiceDef.SetLabel(label)
	return builder
}

var _ ChoiceBuilder = (*chosenBuilder)(nil)

func (builder *chosenBuilder) Unless(guard TransitionGuard) ChoiceBuilder {
	builder.choiceDef.SetUnlessGuard(guard)
	return builder
}

func (builder *chosenBuilder) OnTrue(eventBuilderFn func(eventBuilder EventBuilder)) ChoiceTrueBuilder {
	builder.choiceDef.SetOnTrue(eventBuilderFn)
	return newChoiceTrueBuilder(builder.choiceDef)
}

func (builder *chosenBuilder) OnFalse(eventBuilderFn func(eventBuilder EventBuilder)) ChoiceFalseBuilder {
	builder.choiceDef.SetOnFalse(eventBuilderFn)
	return newChoiceFalseBuilder(builder.choiceDef)
}

// newChoiceTrueBuilder returns a zero-valued instance of
// ChoiceTrueBuilder, which implements ChoiceTrueBuilder.
func newChoiceTrueBuilder(choiceDef *ChoiceDef) ChoiceTrueBuilder {
	return &choiceTrueBuilder{
		choiceDef: choiceDef,
	}
}

// choiceTrueBuilder implements ChoiceTrueBuilder
type choiceTrueBuilder struct {
	choiceDef *ChoiceDef
}

var _ ChoiceTrueBuilder = (*choiceTrueBuilder)(nil)

func (builder *choiceTrueBuilder) OnFalse(eventBuilderFn func(eventBuilder EventBuilder)) {
	builder.choiceDef.SetOnFalse(eventBuilderFn)
}

// newChoiceFalseBuilder returns a zero-valued instance of
// ChoiceFalseBuilder, which implements ChoiceFalseBuilder.
func newChoiceFalseBuilder(choiceDef *ChoiceDef) ChoiceFalseBuilder {
	return &choiceFalseBuilder{
		choiceDef: choiceDef,
	}
}

// choiceFalseBuilder implements ChoiceFalseBuilder
type choiceFalseBuilder struct {
	choiceDef *ChoiceDef
}

var _ ChoiceFalseBuilder = (*choiceFalseBuilder)(nil)

func (builder *choiceFalseBuilder) OnTrue(eventBuilderFn func(eventBuilder EventBuilder)) {
	builder.choiceDef.SetOnTrue(eventBuilderFn)
}
