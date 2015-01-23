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

func Test_VM_GetAll(t *testing.T) {
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

func Test_VM_DestroyEasyVpn_NonExisting(t *testing.T) {
	mockedProvider := test.MockProvider{
		Config: cfg,
		VMs: []provider.VM{
			provider.VM{},
		},
	}

	DestroyEasyVpn(mockedProvider, "does not exist")
}

func Test_VM_WaitForNewVM(t *testing.T) {
	mockedProvider := test.MockProvider{
		Config: cfg,
		VMs: []provider.VM{
			provider.VM{
				Name: "easy-vpn",
				Id:   "mockId",
			},
			provider.VM{
				Name: "mockName",
			},
		},
	}

	var vm provider.VM
	waitForNewVM(mockedProvider, &vm, "easy-vpn")
	if assert.NotNil(t, vm) {
		assert.Equal(t, "mockId", vm.Id)
	}
}

func Test_VM_statusOfVM(t *testing.T) {
	mockedProvider := test.MockProvider{
		Config: cfg,
		VMs: []provider.VM{
			provider.VM{
				Name:   "easy-vpn",
				Id:     "mockId",
				Status: "active",
			},
		},
	}

	vm := provider.VM{
		Id: "mockId",
	}

	statusOfVM(mockedProvider, &vm)
	if assert.NotNil(t, vm) {
		assert.Equal(t, "active", vm.Status)
	}
}
