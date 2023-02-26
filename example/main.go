package main

import (
	"context"

	reload "github.com/ancalabrese/Reload"
	"github.com/hashicorp/go-hclog"
)

type Config struct {
	Disabled bool   `json:"disabled"`
	Port     string `json:"port"`
	Address  string `json:"address"`
	Timeout  int    `json:"opTimeout"`
}

var config *Config

func main() {
	ctx := context.Background()
	config = &Config{}
	l := hclog.Default()
	rc, _ := reload.New(ctx)

	rc.AddConfiguration("config.json", config)

	l.Info("Update any value in ./config.json to receive new configurations")

	for {
		select {
		case err := <-rc.GetErrChannel():
			l.Error("Received", "err", err)
		case conf := <-rc.GetRoloadChan():
			l.Info("Received", "config", conf.FilePath, " updated:", conf.Config)
		}
	}
}
