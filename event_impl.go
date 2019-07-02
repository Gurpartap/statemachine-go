package statemachine

type eventImpl struct {
	name string
	def *EventDef
}

func newEventImpl() *eventImpl {
	return &eventImpl{}
}

// Event implements Event.
func (m *eventImpl) Event() string {
	return m.name
}

// SetEventDef implements MachineBuildable.
func (m *eventImpl) SetEventDef(event string, def *EventDef) {
	m.def = def
}
