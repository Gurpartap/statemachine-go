package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Gurpartap/statemachine-go"
)

type Process struct {
	statemachine.Machine

	ShouldAutoStart  bool
	IsProcessRunning bool
	SkipTicks        bool
}

func NewProcess() *Process {
	process := &Process{}

	process.Machine = statemachine.BuildNewMachine(func(m statemachine.MachineBuilder) {
		m.States("unmonitored", "stopped", "starting", "running", "stopping", "restarting")
		m.InitialState("unmonitored")

		m.BeforeTransition().To("starting").Do(process.SetAutoStartOn).Label("setAutoStartOn()")
		m.BeforeTransition().To("stopping").Do(process.SetAutoStartOff).Label("setAutoStartOff()")
		m.BeforeTransition().To("restarting").Do(process.SetAutoStartOn).Label("setAutoStartOn()")
		m.BeforeTransition().To("unmonitored").Do(process.SetAutoStartOff).Label("setAutoStartOff()")

		m.AfterTransition().To("starting").Do(process.Start).Label("start()")
		m.AfterTransition().To("stopping").Do(process.Stop).Label("stop()")
		m.AfterTransition().To("restarting").Do(process.Restart).Label("restart()")

		m.BeforeTransition().ToAny().Do(process.NotifyTriggers)
		m.AroundTransition().ToAny().Do(process.RecordTransition)
		m.AfterFailure().OnAnyEvent().Do(process.LogFailure)

		m.Event("monitor").Transition().From("unmonitored").To("stopped")
		m.Event("start").Transition().From("unmonitored", "stopped").To("starting")
		m.Event("stop").Transition().From("running").To("stopping")
		m.Event("restart").Transition().From("running", "stopped").To("restarting")
		m.Event("unmonitor").Transition().FromAny().To("unmonitored")

		m.Event("tick").
			TimedEvery(1 * time.Second).
			// SkipUntil(process.SkipTick).
			Choice(&process.IsProcessRunning).Label("isRunning").
			Unless(process.SkipTick).
			OnTrue(func(e statemachine.EventBuilder) {
				e.Transition().From("starting").To("running")
				e.Transition().From("restarting").To("running")
				e.Transition().From("stopping").To("running")
				e.Transition().From("stopped").To("running")
			}).
			OnFalse(func(e statemachine.EventBuilder) {
				e.Transition().From("starting").To("stopped")
				e.Transition().From("restarting").To("stopped")
				e.Transition().From("running").To("stopped")
				e.Transition().From("stopping").To("stopped")
				e.Transition().From("stopped").To("starting").
					If(&process.ShouldAutoStart).Label("shouldAutoStart")
			})
	})

	return process
}

func (process *Process) GetIsAutoStartOn() bool {
	return process.ShouldAutoStart
}

func (process *Process) SetAutoStartOn() {
	process.ShouldAutoStart = true
}

func (process *Process) SetAutoStartOff() {
	process.ShouldAutoStart = false
}

func (process *Process) SkipTick() bool {
	return false
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
	if err != statemachine.ErrNoMatchingTransition {
		fmt.Println("😾 LogFailure:", event.Event(), err)
	}
}

func (process *Process) SubRecordTransition(transition statemachine.Transition, next func()) {
	fmt.Printf("Sub: ✅  RecordTransition: from: %s to: %s\n", transition.From(), transition.To())
	next()
}

func (process *Process) SubLogFailure(event statemachine.Event, err error) {
	fmt.Println("Sub: 😾 LogFailure:", event.Event(), err)
}

func main() {
	process := NewProcess()
	process.SetAutoStartOn()
	// process.SkipTicks = true

	// go func() {
	// 	for {
	// 		_ = process.Fire("tick")
	// 		time.Sleep(1 * time.Second)
	// 	}
	// }()

	_ = process.Fire("monitor")

	// _ = process.Send(statemachine.TriggerEvent{
	// 	Event: "monitor",
	// })

	// _ = process.Send(statemachine.OverrideState{
	// 	State: statemachine.StateMap{
	// 		"stopped": nil,
	// 	},
	// })

	stateJSON, _ := json.MarshalIndent(process.GetStateMap(), "", "  ")
	fmt.Printf("%s\n", stateJSON)

	time.AfterFunc(2*time.Second, func() {
		process.Stop()
	})

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
}
