package statemachine

type MachineDef struct {
	States       []string
	InitialState string
	Events       []*EventDef `json:",omitempty"`

	BeforeCallbacks []*TransitionCallbackDef `json:",omitempty"`
	AroundCallbacks []*TransitionCallbackDef `json:",omitempty"`
	AfterCallbacks  []*TransitionCallbackDef `json:",omitempty"`

	FailureCallbacks []*EventCallbackDef `json:",omitempty"`
}

func (def *MachineDef) SetStates(states ...string) {
	def.States = append(def.States, states...)
}

func (def *MachineDef) SetInitialState(state string) {
	def.InitialState = state
}

func (def *MachineDef) AddEvent(eventDef *EventDef) {
	def.Events = append(def.Events, eventDef)
}

func (def *MachineDef) AddBeforeCallback(CallbackDef *TransitionCallbackDef) {
	def.BeforeCallbacks = append(def.BeforeCallbacks, CallbackDef)
}

func (def *MachineDef) AddAroundCallback(CallbackDef *TransitionCallbackDef) {
	def.AroundCallbacks = append(def.AroundCallbacks, CallbackDef)
}

func (def *MachineDef) AddAfterCallback(CallbackDef *TransitionCallbackDef) {
	def.AfterCallbacks = append(def.AfterCallbacks, CallbackDef)
}

func (def *MachineDef) AddFailureCallback(CallbackDef *EventCallbackDef) {
	def.FailureCallbacks = append(def.FailureCallbacks, CallbackDef)
}
