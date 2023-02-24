package reload

import (
	"context"
	"fmt"

	"github.com/fsnotify/fsnotify"
)

type Monitor struct {
	ctx               context.Context
	watcher           *fsnotify.Watcher
	configCache       *ConfigCache
	returnEventChan   chan<- (*ConfigurationFile)
	returnErrChan     chan<- (error)
	eventChan         chan<- (fsnotify.Event)
	errChan           chan<- (error)
	writeEventChannel chan (*WriteEvent)
	writeEventHandler *WriteEventHandler
}

var m *Monitor

// GetMonitorInstance returns an singleton instance of ConfigMonitor
// or an error if fsnotify fails to initialize
func GetMonitorInstance(
	ctx context.Context,
	eventChan chan<- (*ConfigurationFile),
	errChan chan<- (error)) (*Monitor, error) {
	if m == nil {
		w, err := fsnotify.NewWatcher()
		if err != nil {
			return nil, fmt.Errorf("error initializing config monitor: %w", err)
		}

		configManager := GetCacheInstance()
		writeEventChannel := make(chan (*WriteEvent))
		weh := NewWriteEventHandler(ctx, writeEventChannel)

		m = &Monitor{
			ctx:               ctx,
			watcher:           w,
			configCache:       configManager,
			writeEventHandler: weh,
			returnEventChan:   eventChan,
			returnErrChan:     errChan,
			writeEventChannel: writeEventChannel,
			eventChan:         make(chan<- fsnotify.Event),
			errChan:           make(chan<- error),
		}

		go m.monitorUp()
	}

	return m, nil
}

// TrackNew adds the file path to the monitored paths
func (cm *Monitor) TrackNew(path string, config interface{}) error {
	c, err := NewConfigurationFile(path, config)
	if err != nil {
		return err
	}

	err = cm.watcher.Add(path)
	if err != nil {
		return fmt.Errorf("error adding new resource %s to monitor: %w", path, err)
	}

	cm.configCache.Add(c)

	return nil
}

// Untrack removes a path from the monitored files
func (cm *Monitor) Untrack(path string) {
	cm.watcher.Remove(path)
	cm.configCache.Remove(path)
}

// Stop monitoring files and close channels
func (cm *Monitor) Stop() {
	cm.watcher.Close()
	close(cm.eventChan)
	close(cm.errChan)
	m = nil
}

// monitorUp starts listening for events.
// When an event is received it is redirected to the correct event handler
func (cm *Monitor) monitorUp() {
	for {
		select {
		case <-cm.ctx.Done():
			cm.Stop()
			return

		case event := <-cm.watcher.Events:
			if event.Op.Has(fsnotify.Write) {
				writeEvent, _ := NewWriteEvent(event)
				cm.writeEventChannel <- writeEvent
			}

		case path := <-cm.writeEventHandler.GetRelaodChan():
			cm.returnEventChan <- cm.configCache.Get(path)

		case err := <-cm.writeEventHandler.GetErrChan():
			//Send any error back to the caller.
			cm.errChan <- err
		}
	}
}