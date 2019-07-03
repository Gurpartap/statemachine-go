package statemachine_test

import (
	"log"

	"github.com/Gurpartap/statemachine-go"
)

var processStates = []string{
	"unmonitored", "running", "stopped",
	"starting", "stopping", "restarting",
}

type ExampleProcess struct {
	statemachine.Machine

	IsAutoStartOn    bool
	IsProcessRunning bool
}

func (p *ExampleProcess) GetIsAutoStartOn() bool {
	return p.IsAutoStartOn
}

func (p *ExampleProcess) GetIsProcessRunning() bool {
	return p.IsProcessRunning
}

func (ExampleProcess) Start()   {}
func (ExampleProcess) Stop()    {}
func (ExampleProcess) Restart() {}

func (ExampleProcess) NotifyTriggers(transition statemachine.Transition)   {}
func (ExampleProcess) RecordTransition(transition statemachine.Transition) {}

func (ExampleProcess) LogFailure(err error) {
	log.Println(err)
}

func Example_systemProcess() {
	process := &ExampleProcess{}

	process.Machine = statemachine.BuildNewMachine(func(machineBuilder statemachine.MachineBuilder) {
		machineBuilder.States(processStates...)
		machineBuilder.InitialState("unmonitored")

		machineBuilder.Event("tick", func(eventBuilder statemachine.EventBuilder) {
			eventBuilder.Transition().From("starting").To("running").If(process.GetIsProcessRunning)
			eventBuilder.Transition().From("starting").To("stopped").Unless(process.GetIsProcessRunning)

			// The process failed to die after entering the stopping state.
			// Change the state to reflect reality.
			eventBuilder.Transition().From("running").To("stopped").Unless(process.GetIsProcessRunning)

			eventBuilder.Transition().From("stopping").To("running").If(process.GetIsProcessRunning)
			eventBuilder.Transition().From("stopping").To("stopped").Unless(process.GetIsProcessRunning)

			eventBuilder.Transition().From("stopped").To("running").If(process.GetIsProcessRunning)
			eventBuilder.Transition().From("stopped").To("starting").If(func(transition statemachine.Transition) bool {
				return process.GetIsAutoStartOn() && !process.GetIsProcessRunning()
			})

			eventBuilder.Transition().From("restarting").To("running").If(process.GetIsProcessRunning)
			eventBuilder.Transition().From("restarting").To("stopped").Unless(process.GetIsProcessRunning)
		})

		machineBuilder.Event("monitor", func(eventBuilder statemachine.EventBuilder) {
			eventBuilder.Transition().From("unmonitored").To("stopped")
		})

		machineBuilder.Event("start", func(eventBuilder statemachine.EventBuilder) {
			eventBuilder.Transition().From("unmonitored", "stopped").To("starting")
		})

		machineBuilder.Event("stop", func(eventBuilder statemachine.EventBuilder) {
			eventBuilder.Transition().From("running").To("stopping")
		})

		machineBuilder.Event("restart", func(eventBuilder statemachine.EventBuilder) {
			eventBuilder.Transition().From("running", "stopped").To("restarting")
		})

		machineBuilder.Event("unmonitor", func(eventBuilder statemachine.EventBuilder) {
			eventBuilder.Transition().FromAny().To("unmonitored")
		})

		machineBuilder.BeforeTransition().To("starting").Do(func() { process.IsAutoStartOn = true })
		machineBuilder.AfterTransition().To("starting").Do(func() { process.Start() })

		machineBuilder.BeforeTransition().To("stopping").Do(func() { process.IsAutoStartOn = false })
		machineBuilder.AfterTransition().To("stopping").Do(func() { process.Stop() })

		machineBuilder.BeforeTransition().To("restarting").Do(func() { process.IsAutoStartOn = true })
		machineBuilder.AfterTransition().To("restarting").Do(func() { process.Restart() })

		machineBuilder.BeforeTransition().To("unmonitored").Do(func() { process.IsAutoStartOn = false })

		machineBuilder.BeforeTransition().ToAny().Do(process.NotifyTriggers)
		machineBuilder.AfterTransition().ToAny().Do(process.RecordTransition)

		machineBuilder.AfterFailure().OnAnyEvent().Do(process.LogFailure)
	})
}
