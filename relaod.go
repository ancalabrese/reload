package reload

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"
)

type ReloadConfig struct {
	ctx           context.Context
	logger        hclog.Logger
	configMonitor *monitor
	errChan       chan (error)
	configChan    chan (*ConfigurationFile)
}

// New creates a new reload config starts obsevering for config changes.
// ctx is the scope used for Reload. When ctx is cancelled Reload will stop monitoring and reloading configurations
func New(ctx context.Context) (*ReloadConfig, error) {

	l := hclog.Default()
	l.SetLevel(hclog.Debug)

	configMonitor, err :=
		newMonitor(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize config monitor: %w", err)
	}

	cf := &ReloadConfig{
		logger:        l,
		ctx:           ctx,
		configMonitor: configMonitor,
		errChan:       make(chan error),
		configChan:    make(chan *ConfigurationFile),
	}

	go cf.startEventListner()

	return cf, nil
}

// AddConfiguration adds a new config file to the monitor.
// path is the file path
// config is a json tagged struct where the config file will be marshalled into
func (rc *ReloadConfig) AddConfiguration(path string, config any) error {
	return rc.configMonitor.trackNew(path, config)
}

func (rc *ReloadConfig) GetErrChannel() <-chan (error) {
	return rc.errChan
}

func (rc *ReloadConfig) GetReloadChan() <-chan (*ConfigurationFile) {
	return rc.configChan
}

// Stop will stop the monitor and clean up resources
func (rc *ReloadConfig) Stop() {
	rc.configMonitor.stop()
	close(rc.errChan)
	close(rc.configChan)
}

func (rc *ReloadConfig) startEventListner() {
	for {
		select {
		case <-rc.ctx.Done():
			rc.Stop()
		case err := <-rc.configMonitor.getNewConfigurationError():
			rc.errChan <- err
		case config := <-rc.configMonitor.getNewConfiguration():
			rc.configChan <- config
		}
	}
}
