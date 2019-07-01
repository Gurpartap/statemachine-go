package statemachine_test

import (
	"fmt"

	"github.com/Gurpartap/statemachine-go"
)

func ExampleBuildNewMachine() {
	p := &ExampleProcess{}

	p.Machine = statemachine.BuildNewMachine(func(m statemachine.MachineBuilder) {
		m.States(processStates...)
		m.InitialState("unmonitored")

		// ...
	})

	fmt.Println(p.Machine.GetState())
	// Output: unmonitored
}

func ExampleNewMachine() {
	p := &ExampleProcess{}

	p.Machine = statemachine.NewMachine()
	p.Machine.Build(func(m statemachine.MachineBuilder) {
		m.States(processStates...)
		m.InitialState("unmonitored")

		// ...
	})

	fmt.Println(p.Machine.GetState())
	// Output: unmonitored
}
