package statemachine_test

import (
	"github.com/Gurpartap/statemachine-go"
)

type Process struct {
	statemachine.Machine
}

func ExampleTurnstile() {
	turnstile := &Process{}
	turnstile.Machine = statemachine.BuildNewMachine(func(m statemachine.MachineBuilder) {
		m.InitialState("locked")

		m.States("locked", "unlocked")

		m.Event("insert_coin", func(e statemachine.EventBuilder) {
			e.Transition().From("locked").To("unlocked")
		})

		m.Event("push", func(e statemachine.EventBuilder) {
			e.Transition().From("unlocked").To("locked")
		})
	})
}
