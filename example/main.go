package main

import (
	"context"
	"log"

	"github.com/ancalabrese/reload"
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

type Config3 struct {
	Server struct {
		KeepAlive int    `yaml:"keepaliveperiodseconds"`
		Addr      string `yaml:"listenaddr"`
		Port      int    `yaml:"port"`
	} `yaml:"server"`
}

var config *Config

func main() {

	ctx := context.Background()
	config = &Config{}
	config2 := &Config2{}
	config3 := &Config3{}

	rc, err := reload.New(ctx)
	if err != nil {
		log.Fatal("error:", err)
	}

	err = rc.AddConfiguration("./config.json", config)
	err = rc.AddConfiguration("./config2.json", config2)
	err = rc.AddConfiguration("./config.yaml", config3)

	if err != nil {
		panic(err)
	}

	log.Println("Update any value in ./config.json or ./config2.json to" +
		" receive new configurations")

	for {
		select {
		case err := <-rc.GetErrChannel():
			log.Println("Received err: %w", err)
		case conf := <-rc.GetReloadChan():
			log.Println("Received new config [", conf.FilePath, "]:", conf.Config)
		}
	}
}
