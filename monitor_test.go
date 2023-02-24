package reload

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMonitorInstance_noErrors(t *testing.T) {
	mi, err := GetMonitorInstance(
		context.Background(),
		make(chan<- *ConfigurationFile),
		make(chan<- error))

	assert.NoError(t, err, "Error should be nil")

	assert.NotNil(t, mi.watcher, "watcher is nil")
	assert.Empty(
		t,
		mi.configCache.configurations,
		"Config monitor tracking cache is not empty")

	assert.NotNil(t, mi.eventChan, "eventChan is nil")
	assert.NotNil(t, mi.errChan, "errChan is nil")
	assert.NotNil(t, mi.writeEventHandler, "writeEventHandler is nil")

}
