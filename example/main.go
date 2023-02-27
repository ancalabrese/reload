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
type Config2 struct {
	Setting1 bool   `json:"setting1"`
	Setting2 int    `json:"setting2"`
	Setting3 string `json:"setting3"`
}

var config *Config

func main() {

	ctx := context.Background()
	config = &Config{}
	config2 := &Config2{}

	l := hclog.Default()
	l.SetLevel(hclog.Debug)
	rc, _ := reload.New(ctx)

	rc.AddConfiguration("./config.json", config)
	rc.AddConfiguration("./config2.json", config2)

	l.Info("Update any value in ./config.json or ./config2.json to" +
		"receive new configurations")
	
	for {
		select {
		case err := <-rc.GetErrChannel():
			l.Error("Received", "err", err)
		case conf := <-rc.GetRoloadChan():
			if conf != nil {
				panic(conf)
			}
			l.Debug("Received", "config", conf.FilePath, " updated:", conf.Config)
		}
	}
}
