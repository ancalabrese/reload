package internal

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/ancalabrese/reload/data"
	"github.com/ancalabrese/reload/internal/cache"
	"github.com/ancalabrese/reload/internal/handlers"
	"github.com/fsnotify/fsnotify"
)

type Monitor struct {
	watcher          *fsnotify.Watcher
	configCache      *cache.Cache
	eventHandlers    []handlers.EventHandler
	returnConfigChan chan (*data.ConfigurationFile)
	returnErrChan    chan (error)
}

// NewMonitor initiate a new Monitor
func NewMonitor(ctx context.Context) (*Monitor, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("error initializing config monitor: %w", err)
	}

	configChace := cache.GetInstance()
	eventHandlers := []handlers.EventHandler{
		handlers.NewWriteEventHandler(ctx, fsWatcher.Events),
	}

	m := &Monitor{
		watcher:          fsWatcher,
		configCache:      configChace,
		eventHandlers:    eventHandlers,
		returnConfigChan: make(chan *data.ConfigurationFile),
		returnErrChan:    make(chan error),
	}

	go m.monitorUp(ctx)

	return m, nil
}

// TrackNew adds the file path to the monitored paths
func (cm *Monitor) TrackNew(path string, config any) error {
	c, err := data.NewConfigurationFile(path, config)
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

func (m *Monitor) GetNewConfiguration() <-chan (*data.ConfigurationFile) {
	return m.returnConfigChan
}

func (m *Monitor) GetNewConfigurationError() <-chan (error) {
	return m.returnErrChan
}

func (m *Monitor) monitorUp(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			m.Stop()
		case config := <-m.configCache.GetOnReload():
			m.returnConfigChan <- config
		case err := <-m.configCache.GetError():
			m.returnErrChan <- err
		}
	}
}
