package main

import (
	"fmt"
	"time"

	"github.com/Gurpartap/statemachine-go"
)

type Process struct {
	statemachine.Machine

	IsAutoStartOn    bool
	IsProcessRunning bool
}

func NewProcess() *Process {
	process := &Process{}

	process.Machine = statemachine.BuildNewMachine(func(m statemachine.MachineBuilder) {
		m.States(
			"unmonitored", "running", "stopped",
			"starting", "stopping", "restarting",
		)
		m.InitialState("unmonitored")

		m.Event("monitor", func(e statemachine.EventBuilder) { e.Transition().From("unmonitored").To("stopped") })
		m.Event("start", func(e statemachine.EventBuilder) { e.Transition().From("unmonitored", "stopped").To("starting") })
		m.Event("stop", func(e statemachine.EventBuilder) { e.Transition().From("running").To("stopping") })
		m.Event("restart", func(e statemachine.EventBuilder) { e.Transition().From("running", "stopped").To("restarting") })
		m.Event("unmonitor", func(e statemachine.EventBuilder) { e.Transition().FromAny().To("unmonitored") })

		m.Event("tick", func(e statemachine.EventBuilder) {
			e.Transition().From("starting").To("running").If(&process.IsProcessRunning)
			e.Transition().From("starting").To("stopped").Unless(&process.IsProcessRunning)

			// The process failed to die after entering the stopping state.
			// Change the state to reflect reality.
			e.Transition().From("running").To("stopped").Unless(&process.IsProcessRunning)

			e.Transition().From("stopping").To("running").If(&process.IsProcessRunning)
			e.Transition().From("stopping").To("stopped").Unless(&process.IsProcessRunning)

			e.Transition().From("stopped").To("running").If(&process.IsProcessRunning)
			e.Transition().From("stopped").To("starting").If(&process.IsAutoStartOn).AndUnless(&process.IsProcessRunning)

			e.Transition().From("restarting").To("running").If(&process.IsProcessRunning)
			e.Transition().From("restarting").To("stopped").Unless(&process.IsProcessRunning)
		})

		m.BeforeTransition().FromAny().To("starting").Do(func() { process.IsAutoStartOn = true })
		m.BeforeTransition().FromAny().To("stopping").Do(func() { process.IsAutoStartOn = false })
		m.BeforeTransition().FromAny().To("restarting").Do(func() { process.IsAutoStartOn = true })
		m.BeforeTransition().FromAny().To("unmonitored").Do(func() { process.IsAutoStartOn = false })
		m.BeforeTransition().FromAny().ToAny().Do(process.NotifyTriggers)

		m.AroundTransition().FromAny().ToAny().Do(process.RecordTransition)

		m.AfterTransition().FromAny().To("starting").Do(func() { process.Start() })
		m.AfterTransition().FromAny().To("stopping").Do(func() { process.Stop() })
		m.AfterTransition().FromAny().To("restarting").Do(func() { process.Restart() })

		m.AfterFailure().OnAnyEvent().Do(process.LogFailure)
	})

	return process
}

func (process *Process) GetIsAutoStartOn() bool {
	return process.IsAutoStartOn
}

func (process *Process) Start() {
	fmt.Println("Start()")
	process.IsProcessRunning = true
}

func (process *Process) Stop() {
	fmt.Println("Stop()")
	process.IsProcessRunning = false
}

func (process *Process) Restart() {
	fmt.Println("Restart()")
	process.IsProcessRunning = true
}

func (process *Process) NotifyTriggers() {
	fmt.Println("✅  NotifyTriggers")
}

func (process *Process) RecordTransition(transition statemachine.Transition, next func()) {
	fmt.Printf("✅  RecordTransition: from: %s to: %s\n", transition.From(), transition.To())
	next()
}

func (process *Process) LogFailure(event statemachine.Event, err error) {
	fmt.Println("😾 LogFailure:", event.Event(), err)
}

func main() {
	process := NewProcess()

	go func() {
		for {
			process.Fire("tick")
			time.Sleep(1 * time.Second)
		}
	}()

	process.Fire("monitor")

	time.AfterFunc(2*time.Second, func() {
		process.Fire("start")

		time.AfterFunc(3*time.Second, func() {
			process.Fire("stop")
		})
	})

	for {

	}
}
