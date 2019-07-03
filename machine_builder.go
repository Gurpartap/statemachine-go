package statemachine

// MachineBuildable implementation is able to consume the result of building
// features from dsl.MachineBuilder. MachineBuildable is oblivious to Machine or
// it's implementation.
type MachineBuildable interface {
	SetMachineDef(def *MachineDef)
}

// MachineBuilder provides the ability to define all features of the state
// machine, including states, events, event transitions, and transition
// callbacks. MachineBuilder is oblivious to Machine or it's implementation.
type MachineBuilder interface {
	// Build plugs the collected feature definitions into any object
	// that understands them (implements MachineBuildable). Use this method
	// if you're not using Machine.Build() to define the state machine.
	Build(machine MachineBuildable)

	SetEventDef(event string, def *EventDef)

	// States pre-defines the set of known states. This is optional because
	// all known states will be identified from other definitions.
	States(states ...string)

	// InitialState defines the state that the machine initializes with.
	// Initial state must be defined for every state machine.
	InitialState(state string)

	Submachine(state string, submachineBuilderFn func(submachineBuilder MachineBuilder))

	// Event provides the ability to define possible transitions for an event.
	Event(name string, eventBuilderFn ...func(eventBuilder EventBuilder)) EventBuilder

	BeforeTransition() TransitionCallbackBuilder
	AroundTransition() TransitionCallbackBuilder
	AfterTransition() TransitionCallbackBuilder
	AfterFailure() EventCallbackBuilder
}

// NewMachineBuilder returns a zero-valued instance of machineBuilder, which
// implements MachineBuilder.
func NewMachineBuilder() MachineBuilder {
	return &machineBuilder{
		def: NewMachineDef(),
	}
}

// machineBuilder implements MachineBuilder.
type machineBuilder struct {
	def *MachineDef
}

var _ MachineBuilder = (*machineBuilder)(nil)

func (m *machineBuilder) States(states ...string) {
	m.def.SetStates(states...)
}

func (m *machineBuilder) InitialState(state string) {
	m.def.SetInitialState(state)
}

func (m *machineBuilder) Submachine(state string, submachineBuilderFn func(submachineBuilder MachineBuilder)) {
	submachineBuilder := &machineBuilder{def: NewMachineDef()}
	submachineBuilderFn(submachineBuilder)
	m.def.SetSubmachine(state, submachineBuilder.def)
}

func (m *machineBuilder) Event(name string, eventBuilderFuncs ...func(eventBuilder EventBuilder)) EventBuilder {
	// TODO: Restrict public use of .Build(...) on this instance of eventBuilder.
	eventBuilder := NewEventBuilder(name)
	for _, eventBuilderFunc := range eventBuilderFuncs {
		eventBuilderFunc(eventBuilder)
	}
	e := newEventImpl()
	eventBuilder.Build(e)
	m.def.AddEvent(name, e.def)
	return eventBuilder
}

func (m *machineBuilder) BeforeTransition() TransitionCallbackBuilder {
	transitionCallbackDef := &TransitionCallbackDef{validateFor: "BeforeTransition"}
	m.def.AddBeforeCallback(transitionCallbackDef)
	return newTransitionCallbackBuilder(transitionCallbackDef)
}

func (m *machineBuilder) AroundTransition() TransitionCallbackBuilder {
	transitionCallbackDef := &TransitionCallbackDef{validateFor: "AroundTransition"}
	m.def.AddAroundCallback(transitionCallbackDef)
	return newTransitionCallbackBuilder(transitionCallbackDef)
}

func (m *machineBuilder) AfterTransition() TransitionCallbackBuilder {
	transitionCallbackDef := &TransitionCallbackDef{validateFor: "AfterTransition"}
	m.def.AddAfterCallback(transitionCallbackDef)
	return newTransitionCallbackBuilder(transitionCallbackDef)
}

func (m *machineBuilder) AfterFailure() EventCallbackBuilder {
	transitionCallbackDef := &EventCallbackDef{validateFor: "AfterFailure"}
	m.def.AddFailureCallback(transitionCallbackDef)
	return newEventCallbackBuilder(transitionCallbackDef)
}

func (m *machineBuilder) Build(machine MachineBuildable) {
	machine.SetMachineDef(m.def)
}

func (m *machineBuilder) SetEventDef(event string, def *EventDef) {
	m.def.AddEvent(event, def)
}
