package statemachine

// Transition provides methods for accessing useful information about the
// active transition.
type Transition interface {
	From() string
	To() string
}

var _ Transition = (*transitionImpl)(nil)
