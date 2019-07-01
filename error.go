package statemachine

import (
	"errors"
)

var ErrNotInitialized = errors.New("state machine not initialized")
var ErrNoSuchEvent = errors.New("no such event")
var ErrNoMatchingTransition = errors.New("no matching transition")
var ErrTransitionNotAllowed = errors.New("transition not allowed")
