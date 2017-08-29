package statemachine_test

import (
	"fmt"

	"github.com/Gurpartap/statemachine-go"
)

func ExampleBuildNewMachine() {
	p := &ExampleProcess{}

	p.Machine = statemachine.BuildNewMachine(func(builder statemachine.MachineBuilder) {
		builder.InitialState("unmonitored")

		// ...
	})

	fmt.Println(p.Machine.GetState())
	//// Output: unmonitored
}

func ExampleNewMachine() {
	p := &ExampleProcess{}

	p.Machine = statemachine.NewMachine()
	p.Machine.Build(func(builder statemachine.MachineBuilder) {
		builder.InitialState("unmonitored")

		// ...
	})

	fmt.Println(p.Machine.GetState())
	//// Output: unmonitored
}
