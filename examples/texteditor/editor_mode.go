package main

import (
	"encoding/json"
	"fmt"

	"github.com/Gurpartap/statemachine-go"
)

type EditorMode struct {
	statemachine.Machine
}

func NewEditorMode() *EditorMode {
	editorMode := &EditorMode{}

	editorMode.Machine = statemachine.BuildNewMachine(func(m statemachine.MachineBuilder) {
		m.States("plaintext")
		m.InitialState("plaintext")

		m.Submachine("richtext", func(bold statemachine.MachineBuilder) {
			bold.ID("bold")
			bold.States("on", "off")
			bold.InitialState("off")

			bold.Event("toggle", func(toggle statemachine.EventBuilder) {
				toggle.Transition().From("on").To("off")
				toggle.Transition().From("off").To("on")
			})
		})

		m.Submachine("richtext", func(underline statemachine.MachineBuilder) {
			underline.ID("underline")
			underline.States("on", "off")
			underline.InitialState("off")

			underline.Event("toggle", func(toggle statemachine.EventBuilder) {
				toggle.Transition().From("on").To("off")
				toggle.Transition().From("off").To("on")
			})
		})

		m.Submachine("richtext", func(italics statemachine.MachineBuilder) {
			italics.ID("italics")
			italics.States("on", "off")
			italics.InitialState("off")

			italics.Event("toggle", func(toggle statemachine.EventBuilder) {
				toggle.Transition().From("on").To("off")
				toggle.Transition().From("off").To("on")
			})
		})

		m.Submachine("richtext", func(list statemachine.MachineBuilder) {
			// list.ID("list")
			// list.States("none", "bullets", "numbers")
			// list.InitialState("none")
			//
			// list.Event("set_none").Transition().FromAny().To("none")
			// list.Event("set_bullets").Transition().FromAny().To("bullets")
			// list.Event("set_numbers").Transition().FromAny().To("numbers")

			list.ID("list")
			list.States("not_listed")
			list.InitialState("listed")

			list.Submachine("listed", func(bullets statemachine.MachineBuilder) {
				bullets.ID("bullets")
				bullets.States("on", "off")
				bullets.InitialState("off")

				bullets.Event("toggle", func(toggle statemachine.EventBuilder) {
					toggle.Transition().From("on").To("off")
					toggle.Transition().From("off").To("on")
				})
			})

			list.Submachine("listed", func(numbers statemachine.MachineBuilder) {
				numbers.ID("numbers")
				numbers.States("on", "off")
				numbers.InitialState("off")

				numbers.Event("toggle", func(toggle statemachine.EventBuilder) {
					toggle.Transition().From("on").To("off")
					toggle.Transition().From("off").To("on")
				})
			})

			// list.Event("set_none").Transition().FromAny().To("not_listed")
			// list.Event("set_bullets").Transition().FromAny().To("listed.bullets")
			// list.Event("set_numbers").Transition().FromAny().To("listed.numbers")

			list.Event("toggle", func(toggle statemachine.EventBuilder) {
				// toggle.Transition().From("not_listed").To("listed.bullets")
				// toggle.Transition().From("listed.bullets").To("listed.numbers")
				// toggle.Transition().From("listed.numbers").To("not_listed")
			})
		})

		m.Event("toggle", func(toggle statemachine.EventBuilder) {
			toggle.Transition().From("plaintext").To("richtext")
			toggle.Transition().From("richtext").To("plaintext")
		})

		m.AroundTransition().FromAny().ToAny().Do(editorMode.LogTransition)
		m.AfterFailure().OnAnyEvent().Do(editorMode.LogFailure)
	})

	return editorMode
}

func (process *EditorMode) LogTransition(transition statemachine.Transition, next func()) {
	fmt.Printf("✅  LogTransition: from: %s to: %s\n", transition.From(), transition.To())
	next()
}

func (process *EditorMode) LogFailure(event statemachine.Event, err error) {
	if err != statemachine.ErrNoMatchingTransition {
		fmt.Println("😾 LogFailure:", event.Event(), err)
	}
}

func main() {
	editorMode := NewEditorMode()

	// _ = editorMode.Send(statemachine.TriggerEvent{
	// 	Event: "toggle",
	// })

	if err := editorMode.Send(statemachine.OverrideState{
		State: statemachine.StateMap{
			"richtext": statemachine.StateMap{
				"bold":      "on",
				"underline": "off",
				"italics":   "on",
				"list": statemachine.StateMap{
					"listed": statemachine.StateMap{
						"bullets": "on",
						"numbers": "off",
					},
				},
			},
		},
	}); err != nil {
		fmt.Printf("err = %+v\n", err)
	}

	stateJSON, _ := json.MarshalIndent(editorMode.GetStateMap(), "", "  ")
	fmt.Printf("state = %s\n", stateJSON)

	submachine, err := editorMode.Submachine("list", "bullets")
	if err != nil {
		fmt.Printf("err = %+v\n", err)
	}
	stateJSON, _ = json.MarshalIndent(submachine.GetStateMap(), "", "  ")
	fmt.Printf("list.bullets = %s\n", stateJSON)

	if err := submachine.Send(statemachine.OverrideState{
		State: "off",
	}); err != nil {
		fmt.Printf("err = %+v\n", err)
	}
	stateJSON, _ = json.MarshalIndent(submachine.GetStateMap(), "", "  ")
	fmt.Printf("list.bullets = %s\n", stateJSON)

	// time.AfterFunc(2*time.Second, func() {
	// 	_ = editorMode.Fire("toggle")
	//
	// 	time.AfterFunc(2*time.Second, func() {
	//
	// 	})
	// })

	// done := make(chan os.Signal, 1)
	// signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	// <-done
}
