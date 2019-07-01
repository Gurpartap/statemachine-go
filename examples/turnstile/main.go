package main

import (
	"fmt"

	"github.com/Gurpartap/statemachine-go"
)

type Turnstile struct {
	// Embed the state machine into our turnstile.
	statemachine.Machine

	// DisplayMsg stores the text that is displayed to the user.
	// Contains either "Pay" or "Go".
	DisplayMsg string
}

func main() {
	turnstile := &Turnstile{}
	turnstile.DisplayMsg = "Pay"

	// Here we're embedding `statemachine.Machine` into our Turnstile struct.
	// Embedding, although not necessary, is a suitable way to apply the
	// state machine behaviour to the struct itself. Turnstile's DisplayMsg
	// variable stores the text that is displayed to the user, and may store
	// either "Pay" or "Go", depending on the state.
	//
	// Now let's set our state machine's definitions utilizing `MachineBuilder`:

	turnstile.Machine = statemachine.BuildNewMachine(func(m statemachine.MachineBuilder) {
		// States may be pre-defined here.
		m.States("locked", "unlocked")

		// Initial State is required to start the state machine.
		// Setting initial state does not invoke any callbacks.
		m.InitialState("locked")

		// Events and their transition(s).
		m.Event("insertCoin", func(e statemachine.EventBuilder) {
			e.Transition().From("locked").To("unlocked")
		})
		m.Event("push", func(e statemachine.EventBuilder) {
			e.Transition().From("unlocked").To("locked")
		})

		// Transition callbacks.
		m.AfterTransition().From("locked").To("unlocked").Do(func() {
			turnstile.DisplayMsg = "Go"
		})
		m.AfterTransition().From("unlocked").To("locked").Do(func() {
			turnstile.DisplayMsg = "Pay"
		})

		m.AfterTransition().FromAny().ToAny().Do(func() {
			fmt.Printf("%s, %s\n", turnstile.GetState(), turnstile.DisplayMsg)
		})
	})

	// Now that our turnstile is ready, let's take it for a spin:

	_ = turnstile.Fire("insertCoin") // => unlocked, Go
	_ = turnstile.Fire("push")       // => locked, Pay
}
