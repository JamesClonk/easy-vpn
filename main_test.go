package main

import (
	"flag"
	"testing"

	"github.com/codegangsta/cli"
	"github.com/stretchr/testify/assert"
)

func Test_Main_ParseGlobalOptions(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	set.String("config", "fixtures/config_test.toml", "...")
	set.String("provider", "", "")
	set.String("autoconnect", "", "")
	set.String("idletime", "", "")
	set.String("uptime", "", "")

	assert.Nil(t, set.Parse([]string{"--provider", "AWS"}))
	assert.Nil(t, set.Parse([]string{"--autoconnect", "TRUE"}))
	assert.Nil(t, set.Parse([]string{"--idletime", "123"}))
	assert.Nil(t, set.Parse([]string{"--uptime", "777"}))

	c := cli.NewContext(nil, nil, set)
	config := parseGlobalOptions(c)

	if assert.NotNil(t, config) {
		assert.Equal(t, "AWS", config.Provider)
		assert.Equal(t, true, config.Options.Autoconnect)
		assert.Equal(t, 123, config.Options.Idletime)
		assert.Equal(t, 777, config.Options.Uptime)
	}
}
