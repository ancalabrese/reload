package configuration

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ancalabrese/Reload/data"
	"github.com/stretchr/testify/assert"
)

type config1 struct {
	Disabled bool   `json:"disabled"`
	Port     string `json:"port"`
	Address  string `json:"address"`
	Timeout  int    `json:"opTimeout"`
}

type config2 struct {
	Setting1 bool   `json:"setting1"`
	Setting2 string `json:"setting2"`
}

var (
	cache *ConfigCache
	c1    = &config1{
		Disabled: true,
		Port:     "5050",
		Address:  "myaddress",
		Timeout:  99,
	}
	c2 = &config2{
		Setting1: false,
		Setting2: "setting2",
	}
	c1File     *os.File
	c2File     *os.File
	c2FilePath string
	c1FilePath string
)

func TestGetChaceInstance_multipleCalls_returnsSameInstance(t *testing.T) {
	cache = GetCacheInstance()
	assert.NotNil(t, cache, "Cache should not be nil")

	cache2 := GetCacheInstance()
	assert.NotNil(t, cache2, "Cache should not be nil")
	assert.Equal(t, cache2, cache, "Cache should be a singleton instnace")
}

func TestAdd_multipleFiles_tracksFilesSeparately(t *testing.T) {
	cache = GetCacheInstance()
	c1File = createDummyConfigFile("./", c1)
	c2File = createDummyConfigFile("./", c2)
	c1FilePath, _ = filepath.Abs(c1File.Name())
	c2FilePath, _ = filepath.Abs(c2File.Name())

	configFile1, err := data.NewConfigurationFile(c1FilePath, c1)
	assert.Nil(t, err, fmt.Sprintf("%s is a valid path", c1FilePath))

	configFile2, err := data.NewConfigurationFile(c2FilePath, c2)
	assert.Nil(t, err, fmt.Sprintf("%s is a valid path", c2FilePath))

	cache.Add(configFile1)
	cache.Add(configFile2)
	cached1 := cache.Get(c1FilePath)
	cached2 := cache.Get(c2FilePath)
	assert.NotNil(t, cached1, "Config cache returnd nil element")
	assert.NotNil(t, cached2, "Config cache returnd nil element")
	assert.Equal(t, configFile1, cached1,
		"Cached ConfigFile should be equal to inserted item")
	assert.Equal(t, configFile2, cached2,
		"Cached ConfigFile should be equal to inserted item")

	deleteDummyConfigFile(c1FilePath)
	deleteDummyConfigFile(c2FilePath)
}

func TestReload_validConfiguration_noErrors(t *testing.T) {
	cache = GetCacheInstance()
	c1File = createDummyConfigFile("./", c1)
	c1FilePath, _ = filepath.Abs(c1File.Name())

	configFile1, _ := data.NewConfigurationFile(c1FilePath, c1)
	cache.Add(configFile1)

	newConfig := &config1{
		Disabled: false,
		Port:     "8080",
		Address:  "coolAddress",
		Timeout:  9,
	}

	f, _ := os.OpenFile(configFile1.FilePath, os.O_WRONLY, 0777)
	defer f.Close()

	json.NewEncoder(f).Encode(newConfig)

	cache.Reload(f.Name())
	assert.Equal(
		t,
		newConfig,
		cache.configurations[configFile1.FilePath].Config,
		"Cached configuration didn't match updated config version")

	deleteDummyConfigFile(c1FilePath)
}

func TestRealod_invalidConfig_errors(t *testing.T) {
	cache = GetCacheInstance()
	c1File = createDummyConfigFile("./", c1)
	c1FilePath, _ = filepath.Abs(c1File.Name())

	configFile1, _ := data.NewConfigurationFile(c1FilePath, c1)
	cache.Add(configFile1)

	f, _ := os.OpenFile(configFile1.FilePath, os.O_RDWR, 0777)
	defer f.Close()

	var txt string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		txt += strings.ReplaceAll(line, "99", "\"x\"")
	}
	f.Truncate(0)
	f.Seek(0, 0)
	f.WriteString(txt)

	err := cache.Reload(f.Name())
	assert.NotNil(
		t,
		err,
		"Reloading an invalid json configuration should return an error")

	deleteDummyConfigFile(c1FilePath)
}

func createDummyConfigFile(path string, config interface{}) *os.File {
	f, err := os.CreateTemp(path, "*.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(config)
	if err != nil {
		panic(err)
	}
	return f
}

func deleteDummyConfigFile(filePath string) {
	os.Remove(filePath)
}