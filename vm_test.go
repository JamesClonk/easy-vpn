package main

import (
	"flag"
	"testing"

	"github.com/JamesClonk/easy-vpn/provider"
	"github.com/codegangsta/cli"
	"github.com/stretchr/testify/assert"
)

func Test_Main_GetAllVMs(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	set.String("config", "fixtures/config_test.toml", "...")
	c := cli.NewContext(nil, nil, set)

	cfg := parseGlobalOptions(c)
	mockedProvider1 := MockProvider{
		Config: cfg,
		VMs: []provider.VM{
			provider.VM{},
			provider.VM{
				Name: EASYVPN_IDENTIFIER,
				Id:   "mockId",
			},
			provider.VM{
				Name: "mockName",
			},
		},
	}

	machines1 := getAllVMs(mockedProvider1)
	if assert.NotNil(t, machines1) {
		assert.Equal(t, 3, len(machines1))
	}

	mockedProvider2 := MockProvider{Config: cfg}
	machines2 := getAllVMs(mockedProvider2)
	assert.Nil(t, machines2)
}
