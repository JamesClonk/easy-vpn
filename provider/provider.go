package provider

type SshKey struct {
	Id   string
	Name string
	Key  string
}

type API interface {
	GetInstalledSshKeys() ([]SshKey, error)
}
