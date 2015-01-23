package test

import (
	"time"

	"github.com/JamesClonk/easy-vpn/config"
	"github.com/JamesClonk/easy-vpn/provider"

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

func (m MockProvider) CreateVM(name, os, size, region, sshkey string) (string, error) {
	return name + ":" + os + ":" + size + ":" + region + ":" + sshkey, nil
}

func (m MockProvider) StartVM(id string) error {
	return nil
}

func (m MockProvider) DestroyVM(id string) error {
	return nil
}

func (m MockProvider) Sleep() {
	time.Sleep(5 * time.Millisecond)
}
