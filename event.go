package statemachine

type event interface {
	Event() string
	Description() string
}

