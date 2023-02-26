package reload

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
)

type WriteEvent struct {
	WriteEvent *fsnotify.Event
}

// NewWriteEvent creates and new WriteEvent. Return errors is the
// passed fsnotify.Event is not fsnotify.Write
func NewWriteEvent(event fsnotify.Event) (*WriteEvent, error) {
	if !event.Op.Has(fsnotify.Write) {
		return nil, fmt.Errorf("event is not a write event")
	}

	return &WriteEvent{
		WriteEvent: &event,
	}, nil
}
