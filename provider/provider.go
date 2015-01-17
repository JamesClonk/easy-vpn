package provider

type SshKey struct {
	Id   string
	Name string
	Key  string
}

type API interface {
	GetInstalledSshKeys() ([]SshKey, error)
	InstallNewSshKey(name, key string) (string, error)
	UpdateSshKey(id, name, key string) error
}
