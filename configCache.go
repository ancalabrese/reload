package reload

import (
	"fmt"
	"path/filepath"
	"sync"
)

// ConfigCache is the internal cache of monitored files.
type configCache struct {
	configurations map[string]*ConfigurationFile
	onReloadChan   chan (*ConfigurationFile)
	onErrorChan    chan (error)
}

var configManager *configCache
var lock = &sync.Mutex{}

// getCacheInstance get a singleton instance ConfigCache
func getCacheInstance() *configCache {

	if configManager == nil {
		lock.Lock()
		defer lock.Unlock()
		if configManager == nil { // Once locked check instance is still nil
			configManager = &configCache{
				configurations: make(map[string]*ConfigurationFile),
				onReloadChan:   make(chan *ConfigurationFile),
				onErrorChan:    make(chan error),
			}
		}
	}

	return configManager
}

// add new files to ConfigCache.
func (cm *configCache) add(
	configurations ...*ConfigurationFile) {
	for _, c := range configurations {
		if _, ok := cm.configurations[c.FilePath]; !ok {
			cm.configurations[c.FilePath] = c
		}
	}
}

func (cm *configCache) get(path string) *ConfigurationFile {
	if !filepath.IsAbs(path) {
		path, _ = filepath.Abs(path)
	}

	return cm.configurations[path]
}

// Remove removes files from ConfigCache.
func (cm *configCache) remove(path string) {
	if !filepath.IsAbs(path) {
		path, _ = filepath.Abs(path)
	}
	delete(cm.configurations, path)
}

// Reload reads the config file and updates the
// cached configuration files
func (cm *configCache) reload(path string) {
	if !filepath.IsAbs(path) {
		path, _ = filepath.Abs(path)
	}

	err := cm.get(path).loadConfiguration()
	if err != nil {
		err = fmt.Errorf("error loading new config: %w", err)
		cm.onErrorChan <- err
	}

	cm.onReloadChan <- cm.get(path)
}

func (cm *configCache) getOnReload() <-chan (*ConfigurationFile) {
	return cm.onReloadChan
}

func (cm *configCache) getError() <-chan (error) {
	return cm.onErrorChan
}
