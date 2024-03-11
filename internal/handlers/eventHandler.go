package handlers

import (
	"context"

	"github.com/fsnotify/fsnotify"
)

type EventHandler interface {
	// handleEvent receives new events via the channel.
	// It should take process any event that the handler supports.
	handleEvent(eventCh <-chan (fsnotify.Event),
		ctx context.Context, onCancel context.CancelFunc)
}
