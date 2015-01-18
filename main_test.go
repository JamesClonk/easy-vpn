package main

import (
	"flag"
	"os"
	"testing"

	"github.com/JamesClonk/easy-vpn/config"
	"github.com/JamesClonk/easy-vpn/provider"
	"github.com/codegangsta/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProvider struct {
	Config *config.Config
	Keys   []provider.SshKey
	VMs    []provider.VM
	mock.Mock
}

func (m MockProvider) GetProviderName() string {
	return "mock"
}

func (m MockProvider) GetConfig() *config.Config {
	return m.Config
}

func (m MockProvider) GetInstalledSshKeys() ([]provider.SshKey, error) {
	return m.Keys, nil
}

func (m MockProvider) InstallNewSshKey(name, key string) (string, error) {
	return name + ":" + key, nil
}

func (m MockProvider) UpdateSshKey(id, name, key string) (string, error) {
	return id + ":" + name + ":" + key, nil
}

func (m MockProvider) GetAllVMs() ([]provider.VM, error) {
	return m.VMs, nil
}

func (m MockProvider) CreateVM(name, os, size, region string) (string, error) {
	return name + ":" + os + ":" + size + ":" + region, nil
}

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

	assert.Nil(t, set.Parse([]string{"--config", "fixtures/config_test.toml"}))
	assert.Nil(t, set.Parse([]string{"--provider", "aws"}))
	assert.Nil(t, set.Parse([]string{"--api-key", "abcdef1234567890"}))
	assert.Nil(t, set.Parse([]string{"--autoconnect", "TRUE"}))
	assert.Nil(t, set.Parse([]string{"--idletime", "123"}))
	assert.Nil(t, set.Parse([]string{"--uptime", "777"}))

	c := cli.NewContext(nil, nil, set)

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
}

func Test_Main_GetProvider(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	set.String("config", "fixtures/config_test.toml", "...")
	c1 := cli.NewContext(nil, nil, set)

	cfg := parseGlobalOptions(c1)
	if assert.NotNil(t, cfg) {
		assert.Equal(t, "vultr", cfg.Provider)
	}

	p1 := getProvider(c1)
	if assert.NotNil(t, p1) {
		assert.Equal(t, "VULTR", p1.GetProviderName())
	}

	set.String("provider", "", "")
	assert.Nil(t, set.Parse([]string{"--provider", "digitalocean"}))
	c2 := cli.NewContext(nil, nil, set)

	p2 := getProvider(c2)
	if assert.NotNil(t, p2) {
		assert.Equal(t, "DigitalOcean", p2.GetProviderName())
	}
}
