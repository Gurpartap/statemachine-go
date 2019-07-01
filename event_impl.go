package statemachine

type eventImpl struct {
	def *EventDef
}

func newEventImpl() *eventImpl {
	return &eventImpl{}
}

// Event implements Event.
func (m *eventImpl) Event() string {
	return m.def.Name
}

// SetEventDef implements MachineBuildable.
func (m *eventImpl) SetEventDef(def *EventDef) {
	m.def = def
}

// simpleEvent is used in AfterFailure callbacks.
type simpleEvent struct {
	name string
}

// Event implements Event.
func (s *simpleEvent) Event() string {
	return s.name
}
