package statemachine

type Message interface {
	Value() interface{}
}

var _ Message = TriggerEvent{}
var _ Message = OverrideState{}

type TriggerEvent struct {
	Event string
}

func (e TriggerEvent) Value() interface{} {
	return e.Event
}

type OverrideState struct {
	State interface{}
}

func (e OverrideState) Value() interface{} {
	return e.State
}
