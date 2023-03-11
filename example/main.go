package main

import (
	"context"
	"log"

	reload "github.com/ancalabrese/Reload"
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

	rc, _ := reload.New(ctx)

	rc.AddConfiguration("example/config.json", config)
	rc.AddConfiguration("example/config2.json", config2)

	log.Println("Update any value in ./config.json or ./config2.json to" +
		" receive new configurations")

	for {
		select {
		case err := <-rc.GetErrChannel():
			log.Println("Received err: %w", err)
		case conf := <-rc.GetRoloadChan():
			log.Println("Received new config [", conf.FilePath, "]:", conf.Config)
		}
	}
}
