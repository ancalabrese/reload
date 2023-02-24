package reload

import (
	"fmt"
)

// ConfigCache is the internal cache of monitored files.
type ConfigCache struct {
	configurations map[string]*ConfigurationFile
}

var configManager *ConfigCache

// GetCacheInstance get a singleton instance ConfigCache
func GetCacheInstance() *ConfigCache {

	if configManager == nil {
		configManager = &ConfigCache{
			configurations: make(map[string]*ConfigurationFile),
		}
	}

	return configManager
}

// Add new files to ConfigCache.
func (cm *ConfigCache) Add(
	configurations ...*ConfigurationFile) {
	for _, c := range configurations {
		if _, ok := cm.configurations[c.FilePath]; !ok {
			cm.configurations[c.FilePath] = c
		}
	}
}

func (cm *ConfigCache) Get(path string) *ConfigurationFile {
	return cm.configurations[path]
}

// Remove removes files from ConfigCache.
func (cm *ConfigCache) Remove(path string) {
	delete(cm.configurations, path)
}

// Reload reads the config file and updates the
// cached configuration files
func (cm *ConfigCache) Reload(path string) error {
	err := cm.Get(path).LoadConfiguration()
	if err != nil {
		return fmt.Errorf("error loading new config: %w", err)
	}
	return nil
}
