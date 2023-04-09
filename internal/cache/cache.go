package cache

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/ancalabrese/reload/data"
)

// ConfigCache is the internal cache of monitored files.
type Cache struct {
	configurations map[string]*data.ConfigurationFile
	onReloadChan   chan (*data.ConfigurationFile)
	onErrorChan    chan (error)
}

var configCache *Cache
var lock = &sync.Mutex{}

// GetInstance get a singleton instance ConfigCache
func GetInstance() *Cache {

	if configCache == nil {
		lock.Lock()
		defer lock.Unlock()
		if configCache == nil { // Once locked check instance is still nil
			configCache = &Cache{
				configurations: make(map[string]*data.ConfigurationFile),
				onReloadChan:   make(chan *data.ConfigurationFile),
				onErrorChan:    make(chan error),
			}
		}
	}

	return configCache
}

// add new files to ConfigCache.
func (cm *Cache) Add(
	configurations ...*data.ConfigurationFile) {
	for _, c := range configurations {
		if _, ok := cm.configurations[c.FilePath]; !ok {
			cm.configurations[c.FilePath] = c
			cm.Reload(c.FilePath)
		}
	}
}

func (cm *Cache) Get(path string) *data.ConfigurationFile {
	if !filepath.IsAbs(path) {
		path, _ = filepath.Abs(path)
	}

	return cm.configurations[path]
}

// Remove removes files from ConfigCache.
func (cm *Cache) Remove(path string) {
	if !filepath.IsAbs(path) {
		path, _ = filepath.Abs(path)
	}
	delete(cm.configurations, path)
}

// Reload reads the config file and updates the
// cached configuration files
func (cm *Cache) Reload(path string) {
	if !filepath.IsAbs(path) {
		path, _ = filepath.Abs(path)
	}

	err := cm.Get(path).LoadConfiguration()
	if err != nil {
		err = fmt.Errorf("error loading new config: %w", err)
		cm.onErrorChan <- err
		// Rollback file to the current working config in cache
		cm.Get(path).SaveConfiguration()
		return
	}

	cm.onReloadChan <- cm.Get(path)
}

func (cm *Cache) GetOnReload() <-chan (*data.ConfigurationFile) {
	return cm.onReloadChan
}

func (cm *Cache) GetError() <-chan (error) {
	return cm.onErrorChan
}
