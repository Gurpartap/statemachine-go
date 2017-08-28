package statemachine

type state interface {
	State() string
	Description() string
}

