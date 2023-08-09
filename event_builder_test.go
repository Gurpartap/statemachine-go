package statemachine_test

import (
	"fmt"

	"github.com/Gurpartap/statemachine-go"
)

func ExampleNewEventBuilder() {
	p := &ExampleProcess{}
	p.Machine = statemachine.NewMachine()

	machineBuilder := statemachine.NewMachineBuilder()
	machineBuilder.States(processStates...)
	machineBuilder.InitialState("unmonitored")

	eventBuilder := statemachine.NewEventBuilder("monitor")
	eventBuilder.Transition().From("unmonitored").To("stopped")
	eventBuilder.Build(machineBuilder)

	machineBuilder.Build(p.Machine)

	fmt.Println(p.Machine.GetState())

	if err := p.Machine.Fire("monitor"); err != nil {
		fmt.Println(err)
	}

	fmt.Println(p.Machine.GetState())

	if err := p.Machine.Fire("monitor"); err != nil {
		fmt.Println(err)
	}

	// Output: unmonitored
	// stopped
	// no transition from state=stopped for event=monitor: no matching transition
}
