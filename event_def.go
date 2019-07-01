package statemachine

type EventDef struct {
	Name        string
	Transitions []*TransitionDef
}

func (def *EventDef) AddTransition(transitionDef *TransitionDef) {
	def.Transitions = append(def.Transitions, transitionDef)
}
