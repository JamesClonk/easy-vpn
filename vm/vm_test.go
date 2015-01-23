package vm

import (
	"log"
	"testing"

	"github.com/JamesClonk/easy-vpn/config"
	"github.com/JamesClonk/easy-vpn/provider"
	"github.com/JamesClonk/easy-vpn/test"
	"github.com/stretchr/testify/assert"
)

var cfg *config.Config

func init() {
	var err error
	cfg, err = config.LoadConfiguration("../fixtures/config_test.toml")
	if err != nil {
		log.Println(err)
	}
}

func Test_Main_GetAll(t *testing.T) {
	mockedProvider1 := test.MockProvider{
		Config: cfg,
		VMs: []provider.VM{
			provider.VM{},
			provider.VM{
				Name: "easy-vpn",
				Id:   "mockId",
			},
			provider.VM{
				Name: "mockName",
			},
		},
	}

	machines1 := GetAll(mockedProvider1)
	if assert.NotNil(t, machines1) {
		assert.Equal(t, 3, len(machines1))
	}

	mockedProvider2 := test.MockProvider{Config: cfg}
	machines2 := GetAll(mockedProvider2)
	assert.Nil(t, machines2)
}
