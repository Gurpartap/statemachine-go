package statemachine

// machineBuildable implementation is able to consume the result of building
// features from MachineBuilder. machineBuildable is oblivious to Machine or
// it's implementation.
type machineBuildable interface{}

// MachineBuilder provides the ability to define all features of the state
// machine, including states, events, event transitions, and transition
// callbacks. MachineBuilder is oblivious to Machine or it's implementation.
type MachineBuilder interface {
	// State pre-defines the set of known states. This is optional because
	// all known states will be identified from other definitions.
	State(states ...string)

	// InitialState defines the state that the machine initializes with.
	// Initial state must be defined for every state machine.
	InitialState(state string)

	// Event provides the ability to define possible transitions for an event.
	Event(name string, eventBuilderFn func(eventBuilder EventBuilder))

	BeforeTransition() TransitionCallbackBuilder
	AroundTransition() TransitionCallbackBuilder
	AfterTransition() TransitionCallbackBuilder
	AfterFailure() TransitionCallbackBuilder

	// Build plugs the collected feature definitions into any object
	// that understands them (implements machineBuildable). Use this method
	// if you're not using Machine.Build() to define the state machine.
	Build(machine machineBuildable)
}

// NewMachineBuilder returns a zero-valued instance of machineBuilder, which
// implements MachineBuilder.
func NewMachineBuilder() MachineBuilder {
	return &machineBuilder{}
}

// machineBuilder implements MachineBuilder.
type machineBuilder struct{}

func (m *machineBuilder) State(states ...string) {
	panic("not implemented")
}

func (m *machineBuilder) InitialState(state string) {
	panic("not implemented")
}

func (m *machineBuilder) Event(name string, eventBuilderFn func(eventBuilder EventBuilder)) {
	// TODO: Restrict public use of .Build(...) on this instance of eventBuilder.
	eventBuilder := NewEventBuilder(name)
	eventBuilderFn(eventBuilder)
	eventBuilder.Build(m)
}

func (m *machineBuilder) BeforeTransition() TransitionCallbackBuilder {
	return newTransitionCallbackBuilder()
}

func (m *machineBuilder) AroundTransition() TransitionCallbackBuilder {
	return newTransitionCallbackBuilder()
}

func (m *machineBuilder) AfterTransition() TransitionCallbackBuilder {
	return newTransitionCallbackBuilder()
}

func (m *machineBuilder) AfterFailure() TransitionCallbackBuilder {
	return newTransitionCallbackBuilder()
}

func (m *machineBuilder) Build(machine machineBuildable) {
	panic("not implemented")
}
