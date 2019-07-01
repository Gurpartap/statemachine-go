package statemachine

// Event provides methods for accessing useful information about the
// active event.
type Event interface {
	Event() string
}

var _ Event = (*eventImpl)(nil)
var _ EventBuildable = (*eventImpl)(nil)

var _ Event = (*simpleEvent)(nil)
