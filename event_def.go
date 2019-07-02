package statemachine

type EventDef struct {
	Transitions []*TransitionDef
	// Name        string
}

func (def *EventDef) AddTransition(transitionDef *TransitionDef) {
	def.Transitions = append(def.Transitions, transitionDef)
}
