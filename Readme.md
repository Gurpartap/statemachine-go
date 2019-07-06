<img alt="StateMachine" src="https://user-images.githubusercontent.com/39792/29720866-d95eda6a-89d8-11e7-907e-6f0edd7f4ee3.png" width="538" height="72"/>

StateMachine supports creating productive State Machines In Go

[![GoDoc](https://godoc.org/github.com/Gurpartap/statemachine-go?status.svg)](https://godoc.org/github.com/Gurpartap/statemachine-go)

<!-- TOC -->

- [Introduction](#introduction)
    - [Further Reading](#further-reading)
    - [Installation](#installation)
- [Usage](#usage)
    - [Project Goals](#project-goals)
    - [States and Initial State](#states-and-initial-state)
    - [Events](#events)
    - [Transitions](#transitions)
    - [Transition Guards (Conditions)](#transition-guards-conditions)
    - [Transition Callbacks](#transition-callbacks)
        - [Before Transition](#before-transition)
        - [Around Transition](#around-transition)
        - [After Transition](#after-transition)
	- [Event Callbacks](#event-callbacks)
        - [After Failure](#after-failure)
    - [Matchers](#matchers)
        - [Event Transition Matchers](#event-transition-matchers)
        - [Transition Callback Matchers](#transition-callback-matchers)
        - [Event Callback Matchers](#event-callback-matchers)
        - [Callback Functions](#callback-functions)
- [About](#about)

<!-- /TOC -->

## Introduction

State machines provide an alternative way of thinking about how we code any
job/process/workflow.

Using a state machine for an object, that reacts to events differently based on
its current state, reduces the amount of boilerplate and duct-taping you have
to introduce to your code.

StateMachine package provides a feature complete implementation of
finite-state machines in Go.

__What Is A Finite-State Machine Even?__

> A finite-state machine (FSM) is an abstract machine that can be in exactly one
> of a finite number of states at any given time. The FSM can change from one
> state to another in response to some external inputs; the change from one
> state to another is called a transition. An FSM is defined by a list of its
> states, its initial state, and the conditions for each transition.
>
> — [Wikipedia](https://en.wikipedia.org/wiki/Finite-state_machine)

### Further Reading

- [Finite-state Machine](https://en.wikipedia.org/wiki/Finite-state_machine) on Wikipedia
- [Coding State Machines in C or C++](https://barrgroup.com/Embedded-Systems/How-To/Coding-State-Machines) by Miro Samek
- [Statecharts](https://statecharts.github.io)
- [state_machines Ruby Gem](https://github.com/state-machines/state_machines)
- [State Charts XML Notation](https://www.w3.org/TR/scxml/)
- [Flying Spaghetti Monster](https://en.wikipedia.org/wiki/Flying_Spaghetti_Monster)

### Installation

Run this in your project directory:

```bash
go get -u https://github.com/Gurpartap/statemachine-go
```

Import StateMachine with this line in your Go code:

```go
import "github.com/Gurpartap/statemachine-go"
```

## Usage

### Project Goals

> A complex system that works is invariably found to have evolved from a simple
> system that worked. A complex system designed from scratch never works and
> cannot be patched up to make it work. You have to start over with a working
> simple system.
>
> – [John Gall (1975)](https://en.wikipedia.org/wiki/John_Gall_(author)#Gall.27s_law)

Performance is a fairly significant factor when considering the use of
a third party package. However, an API that I can actually code and design in
my mind, ahead of using it, is just as important to me.

StateMachine's API design and developer productivity take precedence over
its benchmark numbers (especially when compared to a bare metal switch
statement based state machine implementation, which may not take you as far).

For this, StateMachine provides a DSL using its builder objects. These
builders compute and validate the state definitions, and then inject the
result (states, events, transitions, callbacks, etc.) into the state machine
during its initialization. Subsequently, these builders are free to be
garbage collected.

Moreover, the state _machinery_ is not dependent on these DSL builders. State
machines may also be initialized from directly allocating definition structs,
or even parsing them from JSON, along with pre-registered callback references.

A state machine definition comprises of the following components:

- [States and Initial State](#states-and-initial-state)
- [Events](#events)
- [Transitions](#transitions)
- [Transition Guards](#transition-guards-conditions)
- [Transition Callbacks](#transition-callbacks)

Each of these components are covered below, along with their example usage
code.

Adding a state machine is as simple as embedding statemachine.Machine in
your struct, defining states and events, along with their transitions.

```go
type Process struct {
    statemachine.Machine

    // or

    Machine statemachine.Machine
}

func NewProcess() *Process {
    process := &Process{}

    process.Machine = statemachine.NewMachine()
    process.Machine.Build(func(m statemachine.MachineBuilder) {
        // ...
    })

    // or

    process.Machine = statemachine.BuildNewMachine(func(m statemachine.MachineBuilder) {
        // ...
    })

    return process
}
```

States, events, and transitions are defined using a DSL composed of "builders",
including `statemachine.MachineBuilder` and
`statemachine.EventBuilder`. These builders provide a clean and type-safe DSL
for writing the specification of how the state machine functions.

The subsequent examples are a close port of
[my experience](https://github.com/Gurpartap/cognizant/blob/master/lib/cognizant/process.rb#L28-L87)
with using the [state_machines](https://github.com/state-machines/state_machines)
Ruby gem. StateMachine Go package's DSL is highly inspired from this Ruby gem.

### States and Initial State

Possible states in the state machine may be manually defined, along with the
initial state. However, states are also inferred from event transition
definitions.

Initial state is set during the initialization of the state machine, and is
required to be defined in the builder.

```go
process.Machine.Build(func(m statemachine.MachineBuilder) {
    m.States("unmonitored", "running", "stopped")
    m.States("starting", "stopping", "restarting")

    // Initial state must be defined.
    m.InitialState("unmonitored")
})
```

### Events

Events act as a virtual function which when fired, trigger a state transition.

```go
process.Machine.Build(func(m statemachine.MachineBuilder) {
    m.Event("monitor", ... )

    m.Event("start", ... )

    m.Event("stop", ... )

    m.Event("restart", ... )

    m.Event("unmonitor", ... )

    m.Event("tick", ... )
})
```

### Transitions

Transitions represent the change in state when an event is fired.

Note that `.From(states ...string)` accepts variadic args.

```go
process.Machine.Build(func(m statemachine.MachineBuilder) {
    m.Event("monitor", func(e statemachine.EventBuilder) {
        e.Transition().From("unmonitored").To("stopped")
    })

    m.Event("start", func(e statemachine.EventBuilder) {
        // from either of the defined states
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
        // ...
    })
})
```

### Transition Guards (Conditions)

Transition Guards are conditional callbacks which expect a boolean return
value, implying whether or not the transition in context should occur.

```go
type TransitionGuardFnBuilder interface {
    If(guardFunc ...TransitionGuardFunc)
    Unless(guardFunc ...TransitionGuardFunc)
}
```

Valid TransitionGuardFunc signatures:

```go
*bool
func() bool
func(transition statemachine.Transition) bool
```

```go
// Assuming process.IsProcessRunning is a bool variable, and
// process.GetIsProcessRunning is a func returning a bool value.
m.Event("tick", func(e statemachine.EventBuilder) {
    // If guard
    e.Transition().From("starting").To("running").If(&process.IsProcessRunning)

    // Unless guard
    e.Transition().From("starting").To("stopped").Unless(process.GetIsProcessRunning)

    // ...

    e.Transition().From("stopped").To("starting").If(func(t statemachine.Transition) bool {
        return process.ShouldAutoStart && !process.GetIsProcessRunning()
    })

    // or

    e.Transition().From("stopped").To("starting").
        If(&process.ShouldAutoStart).
        AndUnless(&process.IsProcessRunning)

    // ...
})
```

### Transition Callbacks

Transition Callback methods are called before, around, or after a transition.

- [Before Transition](#before-transition)
- [Around Transition](#around-transition)
- [After Transition](#after-transition)

#### Before Transition

`Before Transition` callbacks do not act as a conditional, and a bool return
value will not impact the transition.

Valid TransitionCallbackFunc signatures:

```go
func()
func(m statemachine.Machine)
func(t statemachine.Transition)
func(m statemachine.Machine, t statemachine.Transition)
```

```go
process.Machine.Build(func(m statemachine.MachineBuilder) {
    // ...

    m.BeforeTransition().FromAny().To("stopping").Do(func() {
        process.ShouldAutoStart = false
    })

    // ...
}
```

#### Around Transition

`Around Transition`'s callback provides a next func as input, which must be
called inside the callback. (TODO: Missing to call the method will trigger a runtime
failure with an appropriately described error.)

Valid TransitionCallbackFunc signatures:

```go
func(next func())
func(m statemachine.Machine, next func())
func(t statemachine.Transition, next func())
func(m statemachine.Machine, t statemachine.Transition, next func())
```

```go
process.Machine.Build(func(m statemachine.MachineBuilder) {
    // ...

    m.
        AroundTransition().
        From("starting", "restarting").
        To("running").
        Do(func(next func()) {
            start := time.Now()

            // it'll trigger a failure if next is not called
            next()

            elapsed = time.Since(start)

            log.Printf("it took %s to [re]start the process.\n", elapsed)
        })

    // ...
})
```

#### After Transition

`After Transition` callback is called when the state has successfully
transitioned.

Valid TransitionCallbackFunc signatures:

```go
func()
func(m statemachine.Machine)
func(t statemachine.Transition)
func(m statemachine.Machine, t statemachine.Transition)
```

```go
process.Machine.Build(func(m statemachine.MachineBuilder) {
    // ...

    // notify system admin
    m.AfterTransition().From("running").ToAny().Do(process.DialHome)

    // log all transitions
    m.
        AfterTransition().
        ToAny().
        Do(func(t statemachine.Transition) {
            log.Printf("State changed from '%s' to '%s'.\n", t.From(), t.To())
        })

    // ...
})
```

### Event Callbacks

There is only one Event Callback method, which is called after an event fails
to transition the state.

- [After Failure](#after-failure)

#### After Failure

`After Failure` callback is called when there's an error with event firing.

Valid TransitionCallbackFunc signatures:

```go
func()
func(err error)
func(m statemachine.Machine, err error)
func(t statemachine.Event, err error)
func(m statemachine.Machine, t statemachine.Event, err error)
```

```go
process.Machine.Build(func(m statemachine.MachineBuilder) {
    // ...

    m.AfterFailure().OnAnyEvent().
        Do(func(e statemachine.Event, err error) {
            log.Printf(
                "could not transition with event='%s' err=%+v\n",
                e.Event(),
                err
            )
        })

    // ...
})
```

### Matchers

#### Event Transition Matchers

These may map from one or more `from` states to exactly one `to` state.

```go
type TransitionBuilder interface {
    From(states ...string) TransitionFromBuilder
    FromAny() TransitionFromBuilder
    FromAnyExcept(states ...string) TransitionFromBuilder
}

type TransitionFromBuilder interface {
    ExceptFrom(states ...string) TransitionExceptFromBuilder
    To(state string) TransitionToBuilder
}

type TransitionExceptFromBuilder interface {
    To(state string) TransitionToBuilder
}
```

__Examples:__

```go
e.Transition().From("first_gear").To("second_gear")

e.Transition().From("first_gear", "second_gear", "third_gear").To("stalled")

allGears := vehicle.GetAllGearStates()
e.Transition().From(allGears...).ExceptFrom("neutral_gear").To("stalled")

e.Transition().FromAny().To("stalled")

e.Transition().FromAnyExcept("neutral_gear").To("stalled")
```

#### Transition Callback Matchers

These may map from one or more `from` states to one or more `to` states.

```go
type TransitionCallbackBuilder interface {
    From(states ...string) TransitionCallbackFromBuilder
    FromAny() TransitionCallbackFromBuilder
    FromAnyExcept(states ...string) TransitionCallbackFromBuilder
}

type TransitionCallbackFromBuilder interface {
    ExceptFrom(states ...string) TransitionCallbackExceptFromBuilder
    To(states ...string) TransitionCallbackToBuilder
    ToSame() TransitionCallbackToBuilder
    ToAny() TransitionCallbackToBuilder
    ToAnyExcept(states ...string) TransitionCallbackToBuilder
}

type TransitionCallbackExceptFromBuilder interface {
    To(states ...string) TransitionCallbackToBuilder
    ToSame() TransitionCallbackToBuilder
    ToAny() TransitionCallbackToBuilder
    ToAnyExcept(states ...string) TransitionCallbackToBuilder
}
```

__Examples:__

```go
m.BeforeTransition().From("idle").ToAny().Do(someFunc)

m.AroundTransition().From("state_x").ToAnyExcept("state_y").Do(someFunc)

m.AfterTransition().ToAny().Do(someFunc)
// ...is same as:
m.AfterTransition().FromAny().ToAny().Do(someFunc)
```

#### Event Callback Matchers

These may match on one or more `events`.

```go
type EventCallbackBuilder interface {
	On(events ...string) EventCallbackOnBuilder
	OnAnyEvent() EventCallbackOnBuilder
	OnAnyEventExcept(events ...string) EventCallbackOnBuilder
}

type EventCallbackOnBuilder interface {
	Do(callbackFunc EventCallbackFunc) EventCallbackOnBuilder
}
```

__Examples:__

```go
m.AfterFailure().OnAnyEventExcept("event_z").Do(someFunc)
```

#### Callback Functions

Any callback function's arguments (and return types) are dynamically set based
on what types are defined (dependency injection). Setting any unavailable arg
or return type will cause a panic during initialization.

For example, if your BeforeTransition() callback does not need access to the
`statemachine.Transition` variable, you may just define the callback with a
blank function signature: `func()`, instead of
`func(t statemachine.Transition)`. Similarly, for an AfterFailure()
callback you can use `func(err error)`, or
`func(e statemachine.Event, err error)`, or even just `func()` .

## About

    Copyright 2017 Gurpartap Singh

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

        http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
