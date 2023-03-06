package reload

import "github.com/fsnotify/fsnotify"

type eventHandler interface {
	// handleEvent receives new events via the channel.
	// It should take process any event that the handler supports.
	handleEvent(eventCh <-chan (fsnotify.Event))
}
