package main

import (
	"flag"
	"os/user"
	"testing"

	"github.com/JamesClonk/easy-vpn/config"
	"github.com/JamesClonk/easy-vpn/provider"
	"github.com/codegangsta/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedProvider struct {
	Config *config.Config
	Keys   []provider.SshKey
	mock.Mock
}

func (m MockedProvider) GetProviderName() string {
	return "mock"
}

func (m MockedProvider) GetConfig() *config.Config {
	return m.Config
}

func (m MockedProvider) GetInstalledSshKeys() ([]provider.SshKey, error) {
	return m.Keys, nil
}

func (m MockedProvider) InstallNewSshKey(name, key string) (string, error) {
	return name + ":" + key, nil
}

func (m MockedProvider) UpdateSshKey(id, name, key string) (string, error) {
	return id + ":" + name + ":" + key, nil
}

func Test_Main_GetEasyVpnSshKeyId(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	set.String("config", "fixtures/config_test.toml", "...")
	c := cli.NewContext(nil, nil, set)

	cfg := parseGlobalOptions(c)
	mockedProvider1 := MockedProvider{
		Config: cfg,
		Keys: []provider.SshKey{
			provider.SshKey{},
			provider.SshKey{
				Name: EASYVPN_IDENTIFIER,
				Id:   "mockId",
			},
			provider.SshKey{
				Name: "mockName",
			},
		},
	}

	keyId1 := getEasyVpnSshKeyId(mockedProvider1)
	if assert.NotNil(t, keyId1) {
		assert.Equal(t, "mockId:easy-vpn:this would be a public key!\n;)\n", keyId1)
	}

	mockedProvider2 := MockedProvider{Config: cfg}
	keyId2 := getEasyVpnSshKeyId(mockedProvider2)
	if assert.NotNil(t, keyId2) {
		assert.Equal(t, "easy-vpn:this would be a public key!\n;)\n", keyId2)
	}
}

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
