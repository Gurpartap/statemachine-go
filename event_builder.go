package statemachine

// EventBuildable implementation is able to consume the result of building
// features from EventBuilder. EventBuildable is oblivious to Event or
// it's implementation.
type EventBuildable interface {
	SetEventDef(event string, def *EventDef)
}

// EventBuilder provides the ability to define an event along with its
// transitions and their guards. EventBuilder is oblivious to Event or it's
// implementation.
type EventBuilder interface {
	Timed() TimedEventBuilder

	// Transition begins the transition builder, accepting states and guards.
	Transition() TransitionBuilder

	// Build plugs the collected feature definitions into any object
	// that understands them (implements dsl.EventBuildable). Use this method
	// if you're not using MachineBuilder.Event() to define the event.
	Build(event EventBuildable)
}

// NewEventBuilder returns a zero-valued instance of eventBuilder, which
// implements EventBuilder.
func NewEventBuilder(name string) EventBuilder {
	return &eventBuilder{
		name: name,
		def:  &EventDef{},
	}
}

// eventBuilder implements EventBuilder.
type eventBuilder struct {
	name string
	def  *EventDef
}

var _ EventBuilder = (*eventBuilder)(nil)

func (e *eventBuilder) Timed() TimedEventBuilder {
	return NewTimedEventBuilder(e.def)
}

func (e *eventBuilder) Transition() TransitionBuilder {
	transitionDef := &TransitionDef{}
	e.def.AddTransition(transitionDef)
	return newTransitionBuilder(transitionDef)
}

func (e *eventBuilder) Build(event EventBuildable) {
	event.SetEventDef(e.name, e.def)
}
