package statemachine_test

import (
	"fmt"

	"github.com/Gurpartap/statemachine-go"
)

func ExampleNewEventBuilder() {
	p := &ExampleProcess{}
	p.Machine = statemachine.NewMachine()

	machineBuilder := statemachine.NewMachineBuilder()
	machineBuilder.InitialState("unmonitored")

	eventBuilder := statemachine.NewEventBuilder("monitor")
	eventBuilder.Transition().From("unmonitored").To("stopped")
	eventBuilder.Build(machineBuilder)

	machineBuilder.Build(p.Machine)

	p.Machine.Fire("monitor")

	fmt.Println(p.Machine.GetState())
	//// Output: stopped
}
