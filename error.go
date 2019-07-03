package statemachine

import (
	"errors"
)

var ErrNoMatchingTransition = errors.New("no matching transition")
var ErrTransitionNotAllowed = errors.New("transition not allowed")
var ErrStateTypeNotSupported = errors.New("state type not supported")
