package statemachine

import (
	"time"
)

type TimedEventBuilder interface {
	Every(duration time.Duration)
}

// NewTimedEventBuilder returns a zero-valued instance of timedEventBuilder, which
// implements TimedEventBuilder.
func NewTimedEventBuilder(eventDef *EventDef) TimedEventBuilder {
	return &timedEventBuilder{
		def: eventDef,
	}
}

// timedEventBuilder implements TimedEventBuilder.
type timedEventBuilder struct {
	def *EventDef
}

var _ TimedEventBuilder = (*timedEventBuilder)(nil)

func (te *timedEventBuilder) Every(duration time.Duration) {
	te.def.SetEvery(duration)
}
