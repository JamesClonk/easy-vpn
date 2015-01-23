package config

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

var cfg *Config

func init() {
	var err error
	cfg, err = LoadConfiguration("../fixtures/config_test.toml")
	if err != nil {
		log.Println(err)
	}
}

func Test_Config_LoadConfiguration(t *testing.T) {
	if assert.NotNil(t, cfg) {
		assert.Equal(t, "vultr", cfg.Provider)
		assert.Equal(t, "../fixtures/vps_rsa", cfg.PrivateKeyFile)
		assert.Equal(t, "../fixtures/vps_rsa.pub", cfg.PublicKeyFile)
		assert.Equal(t, 5, cfg.Sleep)
	}
}

func Test_Config_LoadConfiguration_Options(t *testing.T) {
	if assert.NotNil(t, cfg) {
		assert.Equal(t, 20, cfg.Options.Idletime)
		assert.Equal(t, 300, cfg.Options.Uptime)
		assert.Equal(t, false, cfg.Options.Autoconnect)
	}
}

func Test_Config_LoadConfiguration_Providers(t *testing.T) {
	if assert.NotNil(t, cfg) {
		assert.Equal(t, "abcdefg123xyz", cfg.Providers["digitalocean"].ApiKey)
		assert.Equal(t, "nyc3", cfg.Providers["digitalocean"].Region)
		assert.Equal(t, "1024mb", cfg.Providers["digitalocean"].Size)
		assert.Equal(t, "ubuntu-14-10-i386", cfg.Providers["digitalocean"].OS)

		assert.Equal(t, "xyzabcdefg999", cfg.Providers["vultr"].ApiKey)
		assert.Equal(t, "7", cfg.Providers["vultr"].Region)
		assert.Equal(t, "2", cfg.Providers["vultr"].Size)
		assert.Equal(t, "128", cfg.Providers["vultr"].OS)

		assert.Equal(t, "xyz1234567890", cfg.Providers["aws"].ApiKey)
		assert.Equal(t, "9", cfg.Providers["aws"].Region)
		assert.Equal(t, "7", cfg.Providers["aws"].Size)
		assert.Equal(t, "999", cfg.Providers["aws"].OS)
	}
}

func Test_Config_LoadConfiguration_NoFile(t *testing.T) {
	cfg, err := LoadConfiguration("../fixtures/does_not_exist.toml")
	assert.Nil(t, cfg)
	assert.NotNil(t, err)
}
