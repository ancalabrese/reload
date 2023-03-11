
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)


![logo](/assets/logo.png)

# Reload :arrows_counterclockwise:

Reload automates hot-reloading of configuration files.   
It allows to monitor multiple configuration files (currently supporting JSON) 
and automatically reload them whenever changes are made, 
without the need to restart your application.   
Reload offers error handling and auto-rollback functionality to ensure that 
your application doesn't crash if there's an error in your configuration file. 
Reload is actively being developed and maintained.


## Features :computer:

- Configuration hot-reloading
- Multiple config file monitoring
- Configuration validation

## Introduction :information_source:
- [Introducing Reload: A Golang hot-reload library for your configuration files]()

## Usage/Examples

For a fully working example check [here](https://github.com/ancalabrese/Reload/tree/main/example)
```golang
package main

import (
	"context"

	reload "github.com/ancalabrese/Reload"
)

type Config struct {
	Port     string `json:"port"`
	Address  string `json:"address"`
	Timeout  int    `json:"opTimeout"`
}

func main() {

	ctx := context.Background()
	config := &Config{}

  rc, _ := reload.New(ctx)
	rc.AddConfiguration("example/config.json", config)

  	for {
		select {
		case err := <-rc.GetErrChannel():
			// Handle errors
		case conf := <-rc.GetRoloadChan():
			// Reinitialize applicaiton
		}
	}
}

```

## Authors

- [@ancalabrese](https://calabreseantonio.com)

## Appendix :rocket:

Reload is still in active development and currently in beta. 
We’re working on adding even more features and functionality 
to make managing configuration even easier. 
Any feedback and contributions are welcome!

<a href="https://www.buymeacoffee.com/ancalabrese">
  <img src="https://img.buymeacoffee.com/button-api/?text=Buy me pizza&emoji=🍕&slug=ancalabrese&button_colour=5F7FFF&font_colour=ffffff&font_family=Poppins&outline_colour=000000&coffee_colour=FFDD00" width="150px" height="50px" />
</a>

