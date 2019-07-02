package statemachine

type EventDef struct {
	// Name        string
	Transitions []*TransitionDef `json:",omitempty"`
}

func (def *EventDef) AddTransition(transitionDef *TransitionDef) {
	def.Transitions = append(def.Transitions, transitionDef)
}