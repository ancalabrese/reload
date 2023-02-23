package reload

import (
	"context"
	"fmt"

	"github.com/ancalabrese/Reload/configuration"
	"github.com/ancalabrese/Reload/data"
	"github.com/hashicorp/go-hclog"
)

type ReloadConfig struct {
	ctx              context.Context
	logger           hclog.Logger
	errChannel       chan (error)
	configReloadChan chan (*data.ConfigurationFile)
	configMonitor    *configuration.ConfigMonitor
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

	errorChannel := make(chan error)
	configReloadChan := make(chan *data.ConfigurationFile)

	configMonitor, err :=
		configuration.GetConfigMonitorInstance(ctx, configReloadChan, errorChannel)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize config monitor: %w", err)
	}

	cf := &ReloadConfig{
		logger:           l,
		ctx:              ctx,
		errChannel:       errorChannel,
		configReloadChan: configReloadChan,
		configMonitor:    configMonitor,
	}

	return cf, nil
}

// AddConfiguration adds a new config file to the monitor.
// path is the file path
// config is a json tagged struct where the config file will be marshalled into
func (rc *ReloadConfig) AddConfiguration(path string, config interface{}) {
	rc.configMonitor.TrackNew(path, config)
}

func (rc *ReloadConfig) GetErrChannel() <-chan (error) {
	return rc.errChannel
}

func (rc *ReloadConfig) GetRoloadChan() <-chan (*data.ConfigurationFile) {
	return rc.configReloadChan
}

// Close will stop the monitor and clean up resources
func (rc *ReloadConfig) Close() {
	rc.configMonitor.Stop()
	close(rc.errChannel)
	close(rc.configReloadChan)
}
