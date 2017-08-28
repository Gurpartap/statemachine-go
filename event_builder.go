package statemachine

// eventBuildable implementation is able to consume the result of building
// features from EventBuilder. eventBuildable is oblivious to Event or
// it's implementation.
type eventBuildable interface{}

// EventBuilder provides the ability to define an event along with its
// transitions and their guards. EventBuilder is oblivious to Event or it's
// implementation.
type EventBuilder interface {
	// Transition begins the transition builder, accepting states and guards.
	Transition() TransitionBuilder

	// Build plugs the collected feature definitions into any object
	// that understands them (implements eventBuildable). Use this method
	// if you're not using MachineBuilder.Event() to define the event.
	Build(event eventBuildable)
}

// NewEventBuilder returns a zero-valued instance of eventBuilder, which
// implements EventBuilder.
func NewEventBuilder(name string) EventBuilder {
	return &eventBuilder{
		name: name,
	}
}

// eventBuilder implements EventBuilder.
type eventBuilder struct {
	name string
}

func (e *eventBuilder) Transition() TransitionBuilder {
	return newTransitionBuilder()
}

func (e *eventBuilder) Build(event eventBuildable) {
	panic("not implemented")
}
