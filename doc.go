// Package statemachine provides a convenient implementation of a Finite-state Machine (FSM) in Go. Learn more about FSM on Wikipedia:
//
//  https://en.wikipedia.org/wiki/Finite-state_machine
//
// Let's begin with an object that can suitably adopt state machine behaviour:
//
//  type Turnstile struct {
//      // Embed the state machine into our turnstile.
//      statemachine.Machine
//
//      // DisplayMsg stores the text that is displayed to the user.
//      // Contains either "Pay" or "Go".
//      DisplayMsg string
//  }
//
//  t := &Turnstile{}
//  t.DisplayMsg = "Pay"
//
// Here we're embedding `statemachine.Machine` into our Turnstile struct. Embedding, although not necessary, is a suitable way to apply the state machine behaviour to the struct itself. Turnstile's DisplayMsg variable stores the text that is displayed to the user, and may store either "Pay" or "Go", depending on the state.
//
// Now let's set our state machine's definitions utilizing `MachineBuilder`:
//
//  t.BuildNewMachine(func(m statemachine.MachineBuilder) {
//      // States may be pre-defined here.
//      m.States("locked", "unlocked")
//
//      // Initial State is required to start the state machine.
//      // Setting initial state does not invoke any callbacks.
//      m.InitialState("locked")
//
//      // Events and their transition(s).
//      m.Event("insertCoin", func(e statemachine.EventBuilder) {
//          e.Transition().From("locked").To("unlocked")
//      })
//      m.Event("turn", func(e statemachine.EventBuilder) {
//          e.Transition().From("unlocked").To("locked")
//      })
//
//      // Transition callbacks.
//      m.AfterTransition().From("locked").To("unlocked").Do(func() {
//          t.DisplayMsg = "Go"
//      })
//      m.AfterTransition().From("unlocked").To("locked").Do(func() {
//          t.DisplayMsg = "Pay"
//      })
//  })
//
// Now that our turnstile is ready, let's take it for a spin:
//
//  t.StartMachine()
//  fmt.Println(t.State(), t.DisplayMsg) // => locked, Pay
//
//  err := t.Fire("turn")
//  fmt.Println(err)                     // => no matching transition for the event
//  fmt.Println(t.State(), t.DisplayMsg) // => locked, Pay
//
//  t.Fire("insertCoin")
//  fmt.Println(t.State(), t.DisplayMsg) // => unlocked, Go
//
//  t.Fire("turn")
//  fmt.Println(t.State(), t.DisplayMsg) // => locked, Pay
//
//  t.StopMachine()
package statemachine
