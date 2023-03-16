package reload

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type monitor struct {
	ctx              context.Context
	watcher          *fsnotify.Watcher
	configCache      *configCache
	eventHandlers    []eventHandler
	returnConfigChan chan (*ConfigurationFile)
	returnErrChan    chan (error)
}

// newMonitor initiate a new Monitor
func newMonitor(ctx context.Context) (*monitor, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("error initializing config monitor: %w", err)
	}

	configChace := getCacheInstance()
	eventHandlers := []eventHandler{
		newWriteEventHandler(ctx, fsWatcher.Events),
	}

	m := &monitor{
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

// trackNew adds the file path to the monitored paths
func (cm *monitor) trackNew(path string, config any) error {
	c, err := newConfigurationFile(path, config)
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

	cm.configCache.add(c)
	return nil
}

// Untrack removes a path from the monitored files
func (cm *monitor) untrack(path string) {
	if !filepath.IsAbs(path) {
		path, _ = filepath.Abs(path)
	}

	cm.watcher.Remove(path)
	cm.configCache.remove(path)
}

// Stop monitoring files and close channels
func (cm *monitor) stop() {
	cm.watcher.Close()
	close(cm.returnConfigChan)
	close(cm.returnErrChan)
}

func (m *monitor) getNewConfiguration() <-chan (*ConfigurationFile) {
	return m.returnConfigChan
}

func (m *monitor) getNewConfigurationError() <-chan (error) {
	return m.returnErrChan
}

func (m *monitor) monitorUp() {
	for {
		select {
		case <-m.ctx.Done():
			m.stop()
		case config := <-m.configCache.getOnReload():
			m.returnConfigChan <- config
		case err := <-m.configCache.getError():
			m.returnErrChan <- err
		}
	}
}
