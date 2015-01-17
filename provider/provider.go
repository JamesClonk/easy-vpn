package provider

import "github.com/JamesClonk/easy-vpn/config"

type SshKey struct {
	Id   string
	Name string
	Key  string
}

type API interface {
	GetProviderName() string
	GetConfig() *config.Config
	GetInstalledSshKeys() ([]SshKey, error)
	InstallNewSshKey(name, key string) (string, error)
	UpdateSshKey(id, name, key string) (string, error)
}
