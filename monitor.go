package reload

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type Monitor struct {
	ctx              context.Context
	watcher          *fsnotify.Watcher
	configCache      *ConfigCache
	eventHandlers    []eventHandler
	returnConfigChan chan (*ConfigurationFile)
	returnErrChan    chan (error)
}

// NewMonitor initiate a new Monitor
func NewMonitor(ctx context.Context) (*Monitor, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("error initializing config monitor: %w", err)
	}

	configChace := GetCacheInstance()
	eventHandlers := []eventHandler{
		NewWriteEventHandler(ctx, fsWatcher.Events),
	}

	m := &Monitor{
		ctx:              ctx,
		watcher:          fsWatcher,
		configCache:      configChace,
		eventHandlers:    eventHandlers,
		returnConfigChan: make(chan *ConfigurationFile),
		returnErrChan:    make(chan error),
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
	close(cm.returnConfigChan)
	close(cm.returnErrChan)
}

func (m *Monitor) GetNewConfiguration() <-chan (*ConfigurationFile) {
	return m.returnConfigChan
}

func (m *Monitor) GetNewConfigurationError() <-chan (error) {
	return m.returnErrChan
}

func (m *Monitor) monitorUp() {
	for {
		select {
		case <-m.ctx.Done():
			m.Stop()
		case config := <-m.configCache.GetOnReload():
			m.returnConfigChan <- config
		case err := <-m.configCache.GetError():
			m.returnErrChan <- err
		}
	}
}
