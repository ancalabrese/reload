package reload

import (
	"context"
	"fmt"

	"github.com/ancalabrese/Reload/configuration"
	"github.com/hashicorp/go-hclog"
)

type ReloadConfig struct {
	ctx                  context.Context
	logger               hclog.Logger
	errorChannel         chan<- (error)
	configurationChannel chan<- (interface{})
	configMonitor        *configuration.ConfigMonitor
}

type Event int

const (
	CONFIG_UPDATE Event = 0
)

// New creates a new reload config starts obsevering for config changes.
// ctx is the scope used for Reload. When ctx is cancelled Reload will stop monitoring and reloading configurations
func New(ctx context.Context) (*ReloadConfig, error) {

	l := hclog.Default()
	l.SetLevel(hclog.Debug)

	configMonitor, err := configuration.GetConfigMonitorInstance(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize config monitor: %w", err)
	}

	cf := &ReloadConfig{
		logger:               l,
		ctx:                  ctx,
		errorChannel:         make(chan<- error),
		configurationChannel: make(chan<- interface{}),
		configMonitor:        configMonitor,
	}

	return cf, nil
}

// AddConfiguration adds a new config file to the monitor.
// path is the file path
// config is a json tagged struct where the config file will be marshalled into
func (rc *ReloadConfig) AddConfiguration(path string, config interface{}) {
	rc.configMonitor.TrackNew(path, config)
}

// Close will stop the monitor and clean up resources
func (rc *ReloadConfig) Close() {
	rc.configMonitor.Stop()
	close(rc.errorChannel)
	close(rc.configurationChannel)
}
