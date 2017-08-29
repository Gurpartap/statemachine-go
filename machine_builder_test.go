package statemachine_test

import (
	"fmt"

	"github.com/Gurpartap/statemachine-go"
)

func ExampleNewMachineBuilder() {
	machineBuilder := statemachine.NewMachineBuilder()

	machineBuilder.InitialState("unmonitored")

	machineBuilder.Event("monitor", func(e statemachine.EventBuilder) {
		e.Transition().From("unmonitored").To("stopped")
	})

	machineBuilder.Event("unmonitor", func(e statemachine.EventBuilder) {
		e.Transition().FromAny().To("unmonitored")
	})

	p := &ExampleProcess{}
	p.Machine = statemachine.NewMachine()

	machineBuilder.Build(p.Machine)

	p.Machine.Fire("monitor")

	fmt.Println(p.Machine.GetState())
	//// Output: stopped
}

func ExampleMachineBuilder_State() {
	p := &ExampleProcess{}

	p.Machine = statemachine.BuildNewMachine(func(builder statemachine.MachineBuilder) {
		builder.State(
			"unmonitored", "running", "stopped",
			"starting", "stopping", "restarting",
		)
		builder.InitialState("unmonitored")

		// ...
	})

	fmt.Println(p.Machine.GetState())
	//// Output: unmonitored
}

func ExampleMachineBuilder_Event() {
	p := &ExampleProcess{}

	p.Machine = statemachine.BuildNewMachine(func(m statemachine.MachineBuilder) {
		m.InitialState("unmonitored")

		m.Event("tick", func(e statemachine.EventBuilder) {
			e.Transition().From("starting").To("running").If(p.GetIsProcessRunning)
			e.Transition().From("starting").To("stopped").Unless(p.GetIsProcessRunning)

			e.Transition().From("running").To("stopped").Unless(p.GetIsProcessRunning)

			e.Transition().From("stopping").To("running").If(p.GetIsProcessRunning)
			e.Transition().From("stopping").To("stopped").Unless(p.GetIsProcessRunning)

			e.Transition().From("stopped").To("running").If(p.GetIsProcessRunning)
			e.Transition().From("stopped").To("starting").If(func() bool {
				return p.GetIsAutoStartOn() && !p.GetIsProcessRunning()
			})

			e.Transition().From("restarting").To("running").If(p.GetIsProcessRunning)
			e.Transition().From("restarting").To("stopped").Unless(p.GetIsProcessRunning)
		})

		m.Event("monitor", func(e statemachine.EventBuilder) {
			e.Transition().From("unmonitored").To("stopped")
		})

		m.Event("start", func(e statemachine.EventBuilder) {
			e.Transition().From("unmonitored", "stopped").To("starting")
		})

		m.Event("stop", func(e statemachine.EventBuilder) {
			e.Transition().From("running").To("stopping")
		})

		m.Event("restart", func(e statemachine.EventBuilder) {
			e.Transition().From("running", "stopped").To("restarting")
		})

		m.Event("unmonitor", func(e statemachine.EventBuilder) {
			e.Transition().FromAny().To("unmonitored")
		})

		// ...
	})

	p.Machine.Fire("start")

	fmt.Println(p.Machine.GetState())
	//// Output: starting
}

func ExampleMachineBuilder_BeforeTransition() {
	p := &ExampleProcess{}

	p.Machine = statemachine.BuildNewMachine(func(m statemachine.MachineBuilder) {
		m.InitialState("unmonitored")

		// define events here...

		m.BeforeTransition().FromAny().To("starting").Do(func() { p.IsAutoStartOn = true })
		m.AfterTransition().FromAny().To("starting").Do(func() { p.Start() })

		m.BeforeTransition().FromAny().To("stopping").Do(func() { p.IsAutoStartOn = false })
		m.AfterTransition().FromAny().To("stopping").Do(func() { p.Stop() })

		m.BeforeTransition().FromAny().To("restarting").Do(func() { p.IsAutoStartOn = true })
		m.AfterTransition().FromAny().To("restarting").Do(func() { p.Restart() })

		m.BeforeTransition().FromAny().To("unmonitored").Do(func() { p.IsAutoStartOn = false })

		m.BeforeTransition().FromAny().ToAny().Do(p.NotifyTriggers)
		m.AfterTransition().FromAny().ToAny().Do(p.RecordTransition)
		m.AfterFailure().FromAny().ToAny().Do(p.LogFailure)
	})

	fmt.Println(p.Machine.GetState())
	//// Output: unmonitored
}
