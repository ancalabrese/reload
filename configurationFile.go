package reload

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ancalabrese/reload/internal/encoding"
)

type ConfigurationFile struct {
	FilePath string
	Config   any
	codec    encoding.Codec
}

func newConfigurationFile(
	path string,
	configuration any) (*ConfigurationFile, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	// At this point only files are supported. Make sure the path is not a folder
	if isDirectory(f) {
		return nil, fmt.Errorf("%s is a directory or not supported file type", path)
	}

	codec := encoding.New(filepath.Ext(path))
	if codec == nil {
		return nil, fmt.Errorf("%s file type is not supported", f.Name())
	}

	return &ConfigurationFile{
		FilePath: path,
		Config:   configuration,
		codec:    codec,
	}, nil
}

func (cf *ConfigurationFile) loadConfiguration() error {
	c, err := os.Open(cf.FilePath)
	if err != nil {
		return fmt.Errorf("failed to open config file %s: %w", cf.FilePath, err)
	}
	defer c.Close()

	// err = json.NewDecoder(c).Decode(cf.Config)
	err = cf.codec.Decode(c, cf.Config)
	if err != nil {
		return fmt.Errorf("[loadConfiguration] - failed to marshal new config: %w", err)
	}

	return nil
}

// isDirectory checks whether the provided files is a directory.
// Directories are not supported.
func isDirectory(f *os.File) bool {
	stat, err := f.Stat()
	if err != nil {
		return true //error occurred, assuming this is not a supported file
	}

	return stat.IsDir()

}
