module github.com/ancalabrese/reload

go 1.18

require (
	github.com/fsnotify/fsnotify v1.6.0
	github.com/hashicorp/go-hclog v1.4.0
	github.com/stretchr/testify v1.8.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.0.0-20220908164124-27713097b956 // indirect
)

replace github.com/ancalabrese/reload => ./

replace github.com/ancalabrese/internal => ./internal
