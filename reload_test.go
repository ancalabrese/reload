package reload

import (
	"testing"
)

type TestConfig struct {
	Config1 bool   `json:"config1"`
	Config2 string `json:"config2"`
	Config3 int    `json:"config3"`
}

var (
	testConfig = &TestConfig{}
)

func TestNewConfigurationFile(t *testing.T) {
	// configFile, err := NewConfigurationFile(testFilePath, testConfig)
	
}
