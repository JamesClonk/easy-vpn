package main

import (
	"flag"
	"os"
	"testing"

	"github.com/codegangsta/cli"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// there can only be 1 TestMain for the whole package main.
	// setup/teardown everything here thats needed for all tests.
	os.Exit(m.Run())
}

func Test_Main_ParseGlobalOptions(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	set.String("config", "fixtures/config_test.toml", "...")
	set.String("provider", "", "")
	set.String("api-key", "", "")
	set.String("autoconnect", "", "")
	set.String("idletime", "", "")
	set.String("uptime", "", "")

	assert.Nil(t, set.Parse([]string{"--provider", "AWS"}))
	assert.Nil(t, set.Parse([]string{"--api-key", "abcdef1234567890"}))
	assert.Nil(t, set.Parse([]string{"--autoconnect", "TRUE"}))
	assert.Nil(t, set.Parse([]string{"--idletime", "123"}))
	assert.Nil(t, set.Parse([]string{"--uptime", "777"}))

	c := cli.NewContext(nil, nil, set)

	cfg := parseGlobalOptions(c)
	if assert.NotNil(t, cfg) {
		assert.Equal(t, "AWS", cfg.Provider)
		assert.Equal(t, "abcdef1234567890", cfg.Providers[cfg.Provider].ApiKey)
		assert.Equal(t, true, cfg.Options.Autoconnect)
		assert.Equal(t, 123, cfg.Options.Idletime)
		assert.Equal(t, 777, cfg.Options.Uptime)
	}
}

func Test_Main_GetProvider(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	set.String("config", "fixtures/config_test.toml", "...")
	c := cli.NewContext(nil, nil, set)

	cfg := parseGlobalOptions(c)
	if assert.NotNil(t, cfg) {
		assert.Equal(t, "vultr", cfg.Provider)
	}

	p1 := getProvider(cfg)
	if assert.NotNil(t, p1) {
		assert.Equal(t, "VULTR", p1.GetProviderName())
	}

	cfg.Provider = "digitalocean"
	p2 := getProvider(cfg)
	if assert.NotNil(t, p2) {
		assert.Equal(t, "DigitalOcean", p2.GetProviderName())
	}
}
