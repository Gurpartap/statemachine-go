package statemachine

// Machine provides a public interface to the state machine implementation.
// It provides methods to build and access features of the state machine.
type Machine interface {
	Build(machineBuilderFn func(machineBuilder MachineBuilder))
	SetMachineDef(def *MachineDef)

	GetState() string
	SetCurrentState(state string)
	IsState(state string) bool

	Submachine(state string) (Machine, error)

	Fire(event string) error

	// TODO: ctx.ForceShutdownSubmachines(true), etc.
	// FireContext(, event string) error
}

var _ Machine = (*machineImpl)(nil)
var _ MachineBuildable = (*machineImpl)(nil)
