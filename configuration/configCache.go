package configuration

import "fmt"

// ConfigurationManager is the internal tracker of monitored files.
type ConfigCache struct {
	configurations map[string]*ConfigurationFile
}

var configManager *ConfigCache

// GetCacheInstance get a Singleton instance ConfigurationManager
func GetCacheInstance() *ConfigCache {

	if configManager == nil {
		configManager = &ConfigCache{
			configurations: make(map[string]*ConfigurationFile),
		}
	}

	return configManager
}

// TrackNewConfig adds new files to ConfigurationManger.
func (cm *ConfigCache) Add(
	configurations ...*ConfigurationFile) {
	for _, c := range configurations {
		if _, ok := cm.configurations[c.filePath]; !ok {
			cm.configurations[c.filePath] = c
		}
	}
}

// Remove removes files from ConfigurationManger.
func (cm *ConfigCache) Remove(path string) {
	delete(cm.configurations, path)
}

// Reload reads the config file and updates the
// cached configuration files
func (cm *ConfigCache) Reload(path string) error {
	err := cm.configurations[path].LoadConfiguration()
	if err != nil {
		return fmt.Errorf("error loading new config: %w", err)
	}
	return nil
}
