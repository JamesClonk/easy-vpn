package main

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testConfig *Config

func TestMain(m *testing.M) {
	var err error
	testConfig, err = loadConfiguration("fixtures/config_test.toml")
	if err != nil {
		log.Println(err)
	}
	os.Exit(m.Run())
}

func Test_Config_LoadConfiguration(t *testing.T) {
	if assert.NotNil(t, testConfig) {
		assert.Equal(t, "vultr", testConfig.Provider)
	}
}

func Test_Config_LoadConfiguration_Options(t *testing.T) {
	if assert.NotNil(t, testConfig) {
		assert.Equal(t, 20, testConfig.Options.Idletime)
		assert.Equal(t, 300, testConfig.Options.Uptime)
		assert.Equal(t, false, testConfig.Options.Autoconnect)
	}
}

func Test_Config_LoadConfiguration_Providers(t *testing.T) {
	if assert.NotNil(t, testConfig) {
		assert.Equal(t, "abcdefg123xyz", testConfig.Providers["digitalocean"].ApiKey)
		assert.Equal(t, 3, testConfig.Providers["digitalocean"].Region)
		assert.Equal(t, "xyzabcdefg999", testConfig.Providers["vultr"].ApiKey)
		assert.Equal(t, 7, testConfig.Providers["vultr"].Region)
	}
}

func Test_Config_LoadConfiguration_NoFile(t *testing.T) {
	config, err := loadConfiguration("fixtures/does_not_exist.toml")
	assert.Nil(t, config)
	assert.NotNil(t, err)
}
