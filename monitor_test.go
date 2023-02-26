package reload_test

import (
	"context"
	"os"
	"testing"

	reload "github.com/ancalabrese/Reload"
	"github.com/stretchr/testify/assert"
)

func TestGetMonitorInstance_noErrors(t *testing.T) {
	_, err := reload.NewMonitor(context.Background())

	assert.NoError(t, err, "Error should be nil")
}

func TestMonitor_trackNewValidPath_noError(t *testing.T) {
	m, _ := reload.NewMonitor(context.Background())
	configManager := reload.GetCacheInstance()

	f, _ := os.CreateTemp("./", "*.json")
	defer f.Close()

	err := m.TrackNew(f.Name(), nil)
	assert.Nil(t, err, "TrackNew returned error")

	c := configManager.Get(f.Name())
	assert.NotNil(t, c, "Tracking file cache not updated")

	os.Remove(f.Name())
}

func TestMonitor_trackNewInvalidPath_error(t *testing.T) {
	m, _ := reload.NewMonitor(
		context.Background())

	err := m.TrackNew("./invalid.json", nil)
	assert.NotNil(t, err, "TrackNew returned error")
}

func TestMonitor_untrackFile_fileRemovedFromCache(t *testing.T) {
	m, _ := reload.NewMonitor(
		context.Background())
	configManager := reload.GetCacheInstance()

	f, _ := os.CreateTemp("./", "*.json")
	defer f.Close()

	_ = m.TrackNew(f.Name(), nil)
	c := configManager.Get(f.Name())
	assert.NotNil(t, c, "Tracking file cache not updated")

	m.Untrack(f.Name())
	c = configManager.Get(f.Name())
	assert.Nil(t, c, "Tracking file still in cache after untracking")

	os.Remove(f.Name())
}

func TestMonitor_stopMonitor_noEvents(t *testing.T) {}
