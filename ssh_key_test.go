package main

import (
	"flag"
	"os/user"
	"testing"

	"github.com/codegangsta/cli"
	"github.com/stretchr/testify/assert"
)

// TODO: add test for getEasyVpnSshKeyId()

func Test_Main_ReadPublicKey(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	set.String("config", "fixtures/config_test.toml", "...")
	c := cli.NewContext(nil, nil, set)

	cfg := parseGlobalOptions(c)
	if assert.NotNil(t, cfg) {
		assert.Equal(t, "vultr", cfg.Provider)
	}

	key := readPublicKey(cfg)
	if assert.NotNil(t, key) {
		assert.Equal(t, "this would be a public key!\n;)\n", key)
	}
}

func Test_Main_SanitizeFilename(t *testing.T) {
	filename := sanitizeFilename("~/test/123.txt")
	if assert.NotNil(t, filename) {
		usr, _ := user.Current()
		home := usr.HomeDir
		assert.Equal(t, home+"/test/123.txt", filename)
	}
}
