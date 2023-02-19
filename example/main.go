package main

import (
	"context"
	"time"

	reload "github.com/ancalabrese/Reload"
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
	reloadConfig, _ := reload.New(ctx)

	reloadConfig.AddConfiguration("config.json", config)

	time.Sleep(1000 * time.Second)
}
