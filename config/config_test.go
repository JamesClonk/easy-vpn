package config

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testConfig *Config

func init() {
	var err error
	testConfig, err = LoadConfiguration("../fixtures/config_test.toml")
	if err != nil {
		log.Println(err)
	}
}

func Test_Config_LoadConfiguration(t *testing.T) {
	if assert.NotNil(t, testConfig) {
		assert.Equal(t, "vultr", testConfig.Provider)
		assert.Equal(t, "fixtures/vps_rsa", testConfig.PrivateKeyFile)
		assert.Equal(t, "fixtures/vps_rsa.pub", testConfig.PublicKeyFile)
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
		assert.Equal(t, "nyc3", testConfig.Providers["digitalocean"].Region)
		assert.Equal(t, "1024m", testConfig.Providers["digitalocean"].Size)
		assert.Equal(t, "ubuntu-14-10-i386", testConfig.Providers["digitalocean"].OS)

		assert.Equal(t, "xyzabcdefg999", testConfig.Providers["vultr"].ApiKey)
		assert.Equal(t, "7", testConfig.Providers["vultr"].Region)
		assert.Equal(t, "2", testConfig.Providers["vultr"].Size)
		assert.Equal(t, "128", testConfig.Providers["vultr"].OS)

		assert.Equal(t, "xyz1234567890", testConfig.Providers["aws"].ApiKey)
		assert.Equal(t, "9", testConfig.Providers["aws"].Region)
		assert.Equal(t, "7", testConfig.Providers["aws"].Size)
		assert.Equal(t, "999", testConfig.Providers["aws"].OS)
	}
}

func Test_Config_LoadConfiguration_NoFile(t *testing.T) {
	cfg, err := LoadConfiguration("../fixtures/does_not_exist.toml")
	assert.Nil(t, cfg)
	assert.NotNil(t, err)
}
