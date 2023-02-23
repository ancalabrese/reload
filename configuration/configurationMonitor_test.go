package configuration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfigurationMonitorInstance_noErrors(t *testing.T) {
	mi, err := GetConfigMonitorInstance(context.Background())
	assert.NoError(t, err, "Error should be nil")

	assert.NotNil(t, mi.watcher, "watcher is nil")
	assert.Empty(
		t,
		mi.configManager.configurations,
		"Config monitor tracking cache is not empty")

		assert.NotNil(t, mi.eventChan, "eventChan is nil")
		assert.NotNil(t, mi.errChan, "errChan is nil")
		assert.NotNil(t, mi.writeEventHandler, "writeEventHandler is nil")
				
}
