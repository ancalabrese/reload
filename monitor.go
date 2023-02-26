package reload

import (
	"context"
	"fmt"
	"path/filepath"

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

// NewMonitor initiate a new Monitor
func NewMonitor(
	ctx context.Context,
	eventChan chan<- (*ConfigurationFile),
	errChan chan<- (error)) (*Monitor, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("error initializing config monitor: %w", err)
	}

	configManager := GetCacheInstance()
	writeEventChannel := make(chan (*WriteEvent))
	weh := NewWriteEventHandler(ctx, writeEventChannel)

	m := &Monitor{
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

	return m, nil
}

// TrackNew adds the file path to the monitored paths
func (cm *Monitor) TrackNew(path string, config interface{}) error {
	c, err := NewConfigurationFile(path, config)
	if err != nil {
		return err
	}

	err = cm.watcher.Add(c.FilePath)
	if err != nil {
		return fmt.Errorf(
			"error adding new resource %s to monitor: %w",
			c.FilePath,
			err)
	}

	cm.configCache.Add(c)
	return nil
}

// Untrack removes a path from the monitored files
func (cm *Monitor) Untrack(path string) {
	if !filepath.IsAbs(path) {
		path, _ = filepath.Abs(path)
	}

	cm.watcher.Remove(path)
	cm.configCache.Remove(path)
}

// Stop monitoring files and close channels
func (cm *Monitor) Stop() {
	cm.watcher.Close()
	close(cm.eventChan)
	close(cm.errChan)
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
