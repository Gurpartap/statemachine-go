package statemachine

type MachineDef struct {
	ID           string                   `json:"id,omitempty" hcl:"id" hcle:"omitempty"`
	States       []string                 `hcl:"states"`
	InitialState string                   `hcl:"initial_state"`
	Events       map[string]*EventDef     `json:",omitempty" hcl:"event" hcle:"omitempty"`
	Submachines  map[string][]*MachineDef `json:",omitempty" hcl:"submachine" hcle:"omitempty"`

	BeforeCallbacks []*TransitionCallbackDef `json:",omitempty" hcl:"before_callbacks" hcle:"omitempty"`
	AroundCallbacks []*TransitionCallbackDef `json:",omitempty" hcl:"around_callbacks" hcle:"omitempty"`
	AfterCallbacks  []*TransitionCallbackDef `json:",omitempty" hcl:"after_callbacks" hcle:"omitempty"`

	FailureCallbacks []*EventCallbackDef `json:",omitempty" hcl:"failure_callbacks" hcle:"omitempty"`
}

func NewMachineDef() *MachineDef {
	return &MachineDef{
		Events:      map[string]*EventDef{},
		Submachines: map[string][]*MachineDef{},
	}
}

func (def *MachineDef) SetID(id string) {
	def.ID = id
}

func (def *MachineDef) SetStates(states ...string) {
	def.States = append(def.States, states...)
}

func (def *MachineDef) SetInitialState(state string) {
	def.InitialState = state
}

func (def *MachineDef) SetSubmachine(state string, submachine *MachineDef) {
	def.Submachines[state] = append(def.Submachines[state], submachine)
}

func (def *MachineDef) AddEvent(event string, eventDef *EventDef) {
	def.Events[event] = eventDef
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
