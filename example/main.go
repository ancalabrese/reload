package main

import (
	"context"
	"encoding/xml"
	"log"

	"github.com/ancalabrese/reload"
)

type jsonConfig struct {
	Disabled bool   `json:"disabled"`
	Port     string `json:"port"`
	Address  string `json:"address"`
	Timeout  int    `json:"opTimeout"`
}

type yamlConfig struct {
	Server struct {
		KeepAlive int    `yaml:"keepaliveperiodseconds"`
		Addr      string `yaml:"listenaddr"`
		Port      int    `yaml:"port"`
	} `yaml:"server"`
}

type xmlConfig struct {
	XMLName xml.Name `xml:"app"`
	Address string   `yaml:"address"`
	Port    int      `yaml:"port"`
}

var json *jsonConfig

func main() {

	ctx := context.Background()
	json = &jsonConfig{}
	yaml := &yamlConfig{}
	xml := &xmlConfig{}

	rc, err := reload.New(ctx)
	if err != nil {
		log.Fatal("error:", err)
	}

	log.Println("Update any value in ./config.json or ./config2.json to" +
		" receive new configurations")

	go func() {
		for {
			select {
			case err := <-rc.GetErrChannel():
				log.Println("Received err: %w", err)
			case conf := <-rc.GetReloadChan():
				log.Println("Received new config [", conf.FilePath, "]:", conf.Config)
			}
		}
	}()

	err = rc.AddConfiguration("./example/config.json", json)
	err = rc.AddConfiguration("./example/config.yaml", yaml)
	err = rc.AddConfiguration("./example/config.xml", xml)

	if err != nil {
		panic(err)
	}
	<-ctx.Done()
}
