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
	config := parseGlobalOptions(c)

	if assert.NotNil(t, config) {
		assert.Equal(t, "AWS", config.Provider)
		assert.Equal(t, "abcdef1234567890", config.Providers[config.Provider].ApiKey)
		assert.Equal(t, true, config.Options.Autoconnect)
		assert.Equal(t, 123, config.Options.Idletime)
		assert.Equal(t, 777, config.Options.Uptime)
	}
}

// TODO: add test for getProvider()
