package provider

import "github.com/JamesClonk/easy-vpn/config"

type SshKey struct {
	Id   string
	Name string
	Key  string
}

type VM struct {
	Id     string
	Name   string
	OS     string
	IP     string
	Region string
	Status string
}

type API interface {
	GetProviderName() string
	GetConfig() *config.Config

	// ssh-keys
	GetInstalledSshKeys() ([]SshKey, error)
	InstallNewSshKey(name, key string) (string, error)
	UpdateSshKey(id, name, key string) (string, error)

	// machines
	GetAllVMs() ([]VM, error)
	CreateVM(name, os, size, region, sshkey string) (string, error)
	StartVM(id string) error

	// for request rate limiting
	Sleep()
}
