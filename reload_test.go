package reload_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	reload "github.com/ancalabrese/Reload"
	"github.com/stretchr/testify/assert"
)

type TestConfig1 struct {
	Config1 bool   `json:"config1"`
	Config2 string `json:"config2"`
	Config3 int    `json:"config3"`
}

type TestConfig2 struct {
	Config1 bool   `json:"config1"`
	Config2 string `json:"config2"`
	Config3 int    `json:"config3"`
}

var (
	testConfig1 = &TestConfig1{
		Config1: true,
		Config2: "someValue",
		Config3: 10,
	}
	testConfig2 = &TestConfig2{
		Config1: true,
		Config2: "someOtherVal",
		Config3: 99,
	}

	ctx = context.Background()
)

func TestNewRealodConfiguration_noError(t *testing.T) {
	_, err := reload.New(ctx)
	assert.Nil(t, err, "New returned error")

}

func TestAddConfiguration_validPath_noErrors(t *testing.T) {
	configFile1 := createTempFile(testConfig1)
	defer deleteTestFile(*configFile1)
	reload, _ := reload.New(ctx)

	err := reload.AddConfiguration(configFile1.Name(), testConfig1)
	assert.Nil(t, err, "Add configuration returned error")

}

func TestAddConfiguration_invalidPath_errors(t *testing.T) {
	reload, _ := reload.New(ctx)
	err := reload.AddConfiguration("./some/non/real/file.txt", testConfig2)
	assert.NotNil(t, err, "Add configuration returned error")

}

func TestAddConfiguration_multipleValidPaths_noError(t *testing.T) {
	reload, _ := reload.New(ctx)
	configFile1 := createTempFile(testConfig1)
	defer deleteTestFile(*configFile1)
	configFile2 := createTempFile(testConfig2)
	defer deleteTestFile(*configFile2)

	err := reload.AddConfiguration(configFile1.Name(), testConfig1)
	assert.Nil(t, err, "Add configuration returned error")
	err = reload.AddConfiguration(configFile2.Name(), testConfig2)
	assert.Nil(t, err, "Add configuration returned error")

}

// relaod.AddConfiguration(configFile1.Name(), &TestConfig1{})
func createTempFile(config interface{}) *os.File {
	f, err := os.CreateTemp("./", "*.json")
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

func deleteTestFile(f os.File) {
	os.Remove(f.Name())
}
