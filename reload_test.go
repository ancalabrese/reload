package reload

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestConfig struct {
	Config1 bool   `json:"config1"`
	Config2 string `json:"config2"`
	Config3 int    `json:"config3"`
}

const (
	testFilePath = "~/testConfig.json"
)

var (
	testConfig = &TestConfig{}
)

func TestNewConfigurationFile(t *testing.T) {
	// configFile, err := NewConfigurationFile(testFilePath, testConfig)
	
}
