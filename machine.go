package statemachine

// Machine provides a public interface to the state machine implementation.
// It provides methods to build and access features of the state machine.
type Machine interface {
	Build(machineBuilderFn func(machineBuilder MachineBuilder))
	GetState() string
	IsState(state string) bool
	Fire(event string) error
}

