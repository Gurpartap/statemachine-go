package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Gurpartap/statemachine-go"
)

type Process struct {
	statemachine.Machine

	IsAutoStartOn    bool
	IsProcessRunning bool
	SkipTicks        bool
}

func NewProcess() *Process {
	process := &Process{}

	process.Machine = statemachine.BuildNewMachine(func(m statemachine.MachineBuilder) {
		m.States("unmonitored", "stopped", "starting", "running", "stopping", "restarting")
		m.InitialState("unmonitored")

		m.BeforeTransition().To("starting").Do(process.SetAutoStartOn)
		m.BeforeTransition().To("stopping").Do(process.SetAutoStartOff)
		m.BeforeTransition().To("restarting").Do(process.SetAutoStartOn)
		m.BeforeTransition().To("unmonitored").Do(process.SetAutoStartOff)

		m.AfterTransition().To("starting").Do(process.Start)
		m.AfterTransition().To("stopping").Do(process.Stop)
		m.AfterTransition().To("restarting").Do(process.Restart)

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
			Choice(&process.IsProcessRunning).
			Unless(process.SkipTick).
			OnTrue(func(e statemachine.EventBuilder) {
				e.Transition().From("starting").To("running")
				e.Transition().From("stopping").To("running")
				e.Transition().From("stopped").To("running")
				e.Transition().From("restarting").To("running")
			}).
			OnFalse(func(e statemachine.EventBuilder) {
				e.Transition().From("starting").To("stopped")
				// The process failed to die after entering the stopping state.
				// Change the state to reflect reality.
				e.Transition().From("running").To("stopped")
				e.Transition().From("stopping").To("stopped")
				e.Transition().From("stopped").To("starting").If(&process.IsAutoStartOn)
				e.Transition().From("restarting").To("stopped")
			})
	})

	return process
}

func (process *Process) GetIsAutoStartOn() bool {
	return process.IsAutoStartOn
}

func (process *Process) SetAutoStartOn() {
	process.IsAutoStartOn = true
}

func (process *Process) SetAutoStartOff() {
	process.IsAutoStartOn = false
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
	// process.SkipTicks = true

	// go func() {
	// 	for {
	// 		_ = process.Fire("tick")
	// 		time.Sleep(1 * time.Second)
	// 	}
	// }()

	_ = process.Fire("monitor")

	time.AfterFunc(2*time.Second, func() {
		_ = process.Fire("start")

		time.AfterFunc(2*time.Second, func() {
			process.Stop()
		})
	})

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
}
