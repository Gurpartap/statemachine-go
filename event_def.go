package statemachine

import (
	"time"
)

type EventDef struct {
	// Name        string
	TimedEvery  time.Duration    `json:",omitempty" hcl:"timed_every" hcle:"omitempty"`
	Choice      *ChoiceDef       `json:",omitempty" hcl:"choice" hcle:"omitempty"`
	Transitions []*TransitionDef `json:",omitempty" hcl:"transitions" hcle:"omitempty"`
}

func (def *EventDef) SetEvery(duration time.Duration) {
	def.TimedEvery = duration
}

func (def *EventDef) SetChoice(choiceDef *ChoiceDef) {
	def.Choice = choiceDef
}

func (def *EventDef) AddTransition(transitionDef *TransitionDef) {
	def.Transitions = append(def.Transitions, transitionDef)
}
