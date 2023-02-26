package reload

import "github.com/fsnotify/fsnotify"

type eventHandler interface {
	// getEvents is used to listen for new events
	getEvents(eventCh <-chan (fsnotify.Event))
	// handleEvent is used to take action on events. It should send any error
	// back via an Error channel
	handleEvent(fsnotify.Event)
}
