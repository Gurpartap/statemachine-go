package statemachine

// Machine provides a public interface to the state machine implementation.
// It provides methods to build and access features of the state machine.
type Machine interface {
	Build(machineBuilderFn func(machineBuilder MachineBuilder))
	GetState() string
	IsState(state string) bool
	Fire(event string) error
}

// BuildNewMachine creates a zero-valued instance of machine, and builds it
// using the passed machineBuilderFn arg.
func BuildNewMachine(machineBuilderFn func(machineBuilder MachineBuilder)) Machine {
	machine := NewMachine()
	machine.Build(machineBuilderFn)
	return machine
}

// NewMachine returns a zero-valued instance of machine, which implements
// Machine.
func NewMachine() Machine {
	return &machine{}
}

type machine struct{}

func (m *machine) Build(machineBuilderFn func(machineBuilder MachineBuilder)) {
	machineBuilder := NewMachineBuilder()
	machineBuilderFn(machineBuilder)
	machineBuilder.Build(m)
}

func (m *machine) GetState() string {
	panic("not implemented")
}

func (m *machine) IsState(state string) bool {
	panic("not implemented")
}

func (m *machine) Fire(event string) error {
	panic("not implemented")
}
