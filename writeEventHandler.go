package reload

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type writeEventHandler struct {
	ctx         context.Context
	cancelFunc  context.CancelFunc
	configCache *configCache
}

// newWriteEventHandler creates a new WriteEventHandler and waits for new
// Write events to come through.
func newWriteEventHandler(
	ctx context.Context,
	eventChannel <-chan (fsnotify.Event)) eventHandler {

	context, cancelFunc := context.WithCancel(ctx)
	weh := &writeEventHandler{
		ctx:         context,
		cancelFunc:  cancelFunc,
		configCache: getCacheInstance(),
	}

	go weh.handleEvent(eventChannel)
	return weh
}

// handleEvent handles any fsnotify.Write events.
// Write events might come in bursts, so it listens until no more events
// are received for the same file, then it attempts to reload the config file.
func (weh *writeEventHandler) handleEvent(eventCh <-chan (fsnotify.Event)) {
	// Wait 100ms for new events; each new event resets the timer.
	waitFor := 100 * time.Millisecond
	var mu sync.Mutex
	// Traking separate timers [as path → timer] for different files
	timers := make(map[string]*time.Timer)
	// Callback fired by the timer
	cleanUpTimerFunc := func(path string) {
		mu.Lock()
		delete(timers, path)
		mu.Unlock()
	}

	for {
		select {
		case event := <-eventCh:
			{
				// Reject any event that is not Write event
				if !event.Has(fsnotify.Write) {
					continue
				}
				// Get timer
				mu.Lock()
				t, ok := timers[event.Name]
				mu.Unlock()

				// if no timer yet create one.
				if !ok {
					t = time.AfterFunc(math.MaxInt64, func() {
						defer cleanUpTimerFunc(event.Name)
						weh.configCache.reload(event.Name)
					})
					t.Stop()

					mu.Lock()
					timers[event.Name] = t
					mu.Unlock()
				}

				// Reset the timer for this path, so it will start from 100ms again.
				t.Reset(waitFor)
			}
		case <-weh.ctx.Done():
			return
		}
	}
}
