package handlers

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type WriteEvent struct {
	writeEvent *fsnotify.Event
}

type WriteEventHandler struct {
	EventChannel chan (*WriteEvent)
	ErrChan      chan (error)
}

// NewWriteEvent creates and new WriteEvent. Return errors is the
// passed fsnotify.Event is not fsnotify.Write
func NewWriteEvent(event fsnotify.Event) (*WriteEvent, error) {
	if !event.Op.Has(fsnotify.Write) {
		return nil, fmt.Errorf("event is not Write event")
	}

	return &WriteEvent{
		writeEvent: &event,
	}, nil
}

// NewWriteEventHandler creates a new WriteEvent handler ans starts
// listening for new Write events
func NewWriteEventHandler(eventChannel chan (*WriteEvent)) *WriteEventHandler {
	weh := &WriteEventHandler{
		EventChannel: eventChannel,
	}
	go weh.handleEvents()

	return weh
}

// handleEvents handles the write events. Writes might come in bursts.
// It makes sure to waits until no more events are received for the same file,
// then it sends it back to the channel to handle the new config.
func (weh *WriteEventHandler) handleEvents() {
	// Wait 100ms for new events; each new event resets the timer.
	waitFor := 100 * time.Millisecond
	var mu sync.Mutex
	// Traking separate timers [as path → timer] for different files
	timers := make(map[string]*time.Timer)
	// Callback fired by the timer
	sendEvent := func(we *WriteEvent) {
		weh.EventChannel <- we

		mu.Lock()
		delete(timers, we.writeEvent.Name)
		mu.Unlock()
	}

	for we := range weh.EventChannel {
		// Get timer
		mu.Lock()
		t, ok := timers[we.writeEvent.Name]
		mu.Unlock()

		// if no timer yet create one.
		if !ok {
			t = time.AfterFunc(math.MaxInt64, func() { sendEvent(we) })
			t.Stop()

			mu.Lock()
			timers[we.writeEvent.Name] = t
			mu.Unlock()
		}

		// Reset the timer for this path, so it will start from 100ms again.
		t.Reset(waitFor)
	}
}
