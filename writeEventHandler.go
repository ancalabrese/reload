package reload

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

)

type WriteEventHandler struct {
	ctx            context.Context
	configCache    *ConfigCache
	eventChannel   <-chan (*WriteEvent)
	errChan        chan (error)
	reloadPathChan chan (string)
}

// NewWriteEventHandler creates a new WriteEventHandler ans starts
// listening for new Write events.
func NewWriteEventHandler(
	ctx context.Context,
	eventChannel <-chan (*WriteEvent)) *WriteEventHandler {

	weh := &WriteEventHandler{
		ctx:            ctx,
		configCache:    GetCacheInstance(),
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
	handleEventFunc := func(we *WriteEvent) {
		cleanUpTimerFunc := func(path string) {
			mu.Lock()
			delete(timers, we.WriteEvent.Name)
			mu.Unlock()
		}
		defer cleanUpTimerFunc(we.WriteEvent.Name)

		err := weh.configCache.Reload(we.WriteEvent.Name)
		if err != nil {
			weh.errChan <- fmt.Errorf("event handler error: %w", err)
			return
		}

		weh.reloadPathChan <- we.WriteEvent.Name
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
				t, ok := timers[we.WriteEvent.Name]
				mu.Unlock()

				// if no timer yet create one.
				if !ok {
					t = time.AfterFunc(math.MaxInt64, func() { handleEventFunc(we) })
					t.Stop()

					mu.Lock()
					timers[we.WriteEvent.Name] = t
					mu.Unlock()
				}

				// Reset the timer for this path, so it will start from 100ms again.
				t.Reset(waitFor)
			}
		}
	}
}
