package main

import (
	"flag"
	"testing"

	"github.com/codegangsta/cli"
	"github.com/stretchr/testify/assert"
)

func Test_Main_Connect(t *testing.T) {
	connect([][]string{[]string{"echo", "hello world!"}}, "123.456.789", "testuser", "testpassword")
}

func Test_Main_ReplaceCommandVariables(t *testing.T) {
	result := replaceCommandVariables(
		[][]string{[]string{"connect", ";$IP;", ":$USER:", " $PASS "}, []string{"disconnect"}},
		"123.456.789",
		"testuser",
		"testpassword")
	if assert.NotNil(t, result) {
		assert.Equal(t, 2, len(result))
		assert.Equal(t, "connect", result[0][0])
		assert.Equal(t, ";123.456.789;", result[0][1])
		assert.Equal(t, ":testuser:", result[0][2])
		assert.Equal(t, " testpassword ", result[0][3])
		assert.Equal(t, "disconnect", result[1][0])
	}
}

func Test_Main_ParseGlobalOptions(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	set.String("config", "fixtures/config_test.toml", "...")
	set.String("provider", "", "")
	set.String("api-key", "", "")
	set.String("autoconnect", "", "")
	set.String("idletime", "", "")
	set.String("uptime", "", "")

	assert.Nil(t, set.Parse([]string{"--config", "fixtures/config_test.toml"}))
	assert.Nil(t, set.Parse([]string{"--provider", "aws"}))
	assert.Nil(t, set.Parse([]string{"--api-key", "abcdef1234567890"}))
	assert.Nil(t, set.Parse([]string{"--autoconnect", "TRUE"}))
	assert.Nil(t, set.Parse([]string{"--idletime", "123"}))
	assert.Nil(t, set.Parse([]string{"--uptime", "777"}))

	c := cli.NewContext(nil, set, set)

	cfg := parseGlobalOptions(c)
	if assert.NotNil(t, cfg) {
		assert.Equal(t, "aws", cfg.Provider)
		assert.Equal(t, "abcdef1234567890", cfg.Providers[cfg.Provider].ApiKey)
		assert.Equal(t, "9", cfg.Providers[cfg.Provider].Region)
		assert.Equal(t, "7", cfg.Providers[cfg.Provider].Size)
		assert.Equal(t, "999", cfg.Providers[cfg.Provider].OS)
		assert.Equal(t, true, cfg.Options.Autoconnect)
		assert.Equal(t, 123, cfg.Options.Idletime)
		assert.Equal(t, 777, cfg.Options.Uptime)
	}

	rset := flag.NewFlagSet("test", 0)
	rset.String("region", "", "")
	assert.Nil(t, rset.Parse([]string{"--region", "888"}))

	c = cli.NewContext(nil, rset, set)

	cfg = parseGlobalOptions(c)
	if assert.NotNil(t, cfg) {
		assert.Equal(t, "aws", cfg.Provider)
		assert.Equal(t, "888", cfg.Providers[cfg.Provider].Region)
	}
}

func Test_Main_GetProvider(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	set.String("config", "fixtures/config_test.toml", "...")
	c1 := cli.NewContext(nil, set, set)

	cfg := parseGlobalOptions(c1)
	if assert.NotNil(t, cfg) {
		assert.Equal(t, "vultr", cfg.Provider)
	}

	p1 := getProvider(c1)
	if assert.NotNil(t, p1) {
		assert.Equal(t, "vultr", p1.GetProviderName())
	}

	set.String("provider", "", "")
	assert.Nil(t, set.Parse([]string{"--provider", "digitalocean"}))
	c2 := cli.NewContext(nil, set, set)

	p2 := getProvider(c2)
	if assert.NotNil(t, p2) {
		assert.Equal(t, "digitalocean", p2.GetProviderName())
	}
}
