package statemachine

// Transition provides methods for accessing useful information about the
// active transition.
type Transition interface {
	Transition() string
	Description() string
	From() string
	To() string
	Args() []interface{}
	Machine() Machine
}

