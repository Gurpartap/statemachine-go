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
}

func NewProcess() *Process {
	process := &Process{}

	process.Machine = statemachine.BuildNewMachine(func(m statemachine.MachineBuilder) {
		m.States(
			"unmonitored", "stopped", "starting", "stopping", "restarting",
		)
		m.InitialState("unmonitored")

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
		// m.BeforeTransition().FromAny().ToAny().Do(process.NotifyTriggers)

		m.AroundTransition().FromAny().ToAny().Do(process.RecordTransition)

		m.AfterTransition().FromAny().To("starting").Do(func() { process.Start() })
		m.AfterTransition().FromAny().To("stopping").Do(func() { process.Stop() })
		m.AfterTransition().FromAny().To("restarting").Do(func() { process.Restart() })

		m.AfterTransition().FromAny().To("running").Do(func() {
			// time.AfterFunc(3*time.Second, func() {
			if submachine, _ := process.Submachine("running"); submachine != nil {
				submachine.Fire("process")
			}
			// })
		})

		m.AfterFailure().OnAnyEvent().Do(process.LogFailure)

		m.Submachine("running", func(sm statemachine.MachineBuilder) {
			sm.States("pending", "success", "failure")
			sm.InitialState("pending")

			sm.AfterTransition().FromAny().To("processing").Do(func() {
				fmt.Println("processing...")
				time.AfterFunc(3*time.Second, func() {
					if submachine, _ := process.Submachine("running"); submachine != nil {
						if subsubmachine, _ := submachine.Submachine("processing"); subsubmachine != nil {
							subsubmachine.Fire("subsubprocess")
						}
					}
				})
			})
			sm.Event("process", func(e statemachine.EventBuilder) {
				e.Transition().From("pending").To("processing")
			})
			sm.Event("succeed", func(e statemachine.EventBuilder) {
				e.Transition().From("processing").To("success")
			})
			sm.Event("fail", func(e statemachine.EventBuilder) {
				e.Transition().From("processing").To("failure")
			})
			sm.AroundTransition().FromAny().ToAny().Do(process.SubRecordTransition)
			sm.AfterFailure().OnAnyEvent().Do(process.SubLogFailure)
			sm.AfterTransition().FromAny().To("success").ExitInto("stopped")
			sm.AfterTransition().FromAny().To("failure").ExitInto("retrying")

			sm.Submachine("processing", func(subsub statemachine.MachineBuilder) {
				subsub.States("loading", "subsubprocessing")
				subsub.InitialState("loading")
				subsub.Event("subsubprocess", func(e statemachine.EventBuilder) {
					e.Transition().From("loading").To("subsubprocessing")
				})
				subsub.Event("to_done", func(e statemachine.EventBuilder) {
					e.Transition().From("subsubprocessing").To("done")
				})
				subsub.AfterTransition().FromAny().To("subsubprocessing").Do(func() {
					fmt.Println("subsubprocessing...")
					time.AfterFunc(3*time.Second, func() {
						if submachine, _ := process.Submachine("running"); submachine != nil {
							if subsubmachine, _ := submachine.Submachine("processing"); subsubmachine != nil {
								subsubmachine.Fire("to_done")
							}
						}
					})
				})
				subsub.AroundTransition().FromAny().ToAny().Do(func(transition statemachine.Transition, next func()) {
					fmt.Printf("SubSub: ✅  RecordTransition: from: %s to: %s\n", transition.From(), transition.To())
					next()
				})
				subsub.AfterFailure().OnAnyEvent().Do(func(event statemachine.Event, err error) {
					fmt.Println("SubSub: 😾 LogFailure:", event.Event(), err)
				})
				subsub.AfterTransition().FromAny().To("done").ExitInto("success")
			})
		})
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

	go func() {
		for {
			process.Fire("tick")
			time.Sleep(1 * time.Second)
		}
	}()

	process.Fire("monitor")

	time.AfterFunc(2*time.Second, func() {
		process.Fire("start")
	})

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
}
