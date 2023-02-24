package handlers

import (
	"context"
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
	ctx            context.Context
	eventChannel   chan (*WriteEvent)
	errChan        chan (error)
	reloadPathChan chan (string)
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

// NewWriteEventHandler creates a new WriteEventHandler ans starts
// listening for new Write events.
func NewWriteEventHandler(
	ctx context.Context,
	eventChannel chan (*WriteEvent)) *WriteEventHandler {
	weh := &WriteEventHandler{
		ctx:            ctx,
		eventChannel:   eventChannel,
		reloadPathChan: make(chan string),
	}

	go weh.handleEvents()

	return weh
}

func (weh *WriteEventHandler) GetRelaodChan() <-chan (string) {
	return weh.reloadPathChan
}

func (weh *WriteEventHandler) GetErrChan() <-chan (error) {
	return weh.errChan
}

// handleEvents handles write events. Writes might come in bursts.
// It makes sure to wait until no more events are received for the same file,
// then it notify watchers via the new config channel.
func (weh *WriteEventHandler) handleEvents() {
	// Wait 100ms for new events; each new event resets the timer.
	waitFor := 100 * time.Millisecond
	var mu sync.Mutex
	// Traking separate timers [as path → timer] for different files
	timers := make(map[string]*time.Timer)
	// Callback fired by the timer
	sendEvent := func(we *WriteEvent) {
		weh.reloadPathChan <- we.writeEvent.Name

		mu.Lock()
		delete(timers, we.writeEvent.Name)
		mu.Unlock()
	}

	for {
		select {
		case <-weh.ctx.Done():
			close(weh.reloadPathChan)
			return

		case we := <-weh.eventChannel:
			{
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
	}
}
