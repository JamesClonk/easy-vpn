package digitalocean

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/JamesClonk/easy-vpn/config"
	"github.com/JamesClonk/easy-vpn/provider"
)

type SshKeys struct {
	Keys []SshKey `json:"ssh_keys"`
}

type SshKey struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Key         string `json:"public_key"`
	Fingerprint string `json:"fingerprint"`
}

type DO struct {
	Config *config.Config
}

func (d DO) GetProviderName() string {
	return "DigitalOcean"
}

func (d DO) GetInstalledSshKeys() (data []provider.SshKey, err error) {
	resp, err := d.doGet(`https://api.digitalocean.com/v2/account/keys`)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}

	var doKeys SshKeys
	if err := json.Unmarshal(body, &doKeys); err != nil {
		return nil, err
	}

	// convert digitalocean ssh-keys into array of provider api ssh-keys
	for _, value := range doKeys.Keys {
		key := provider.SshKey{
			Id:   fmt.Sprintf("%d", value.Id),
			Name: value.Name,
			Key:  value.Key,
		}
		data = append(data, key)
	}

	return data, nil
}

func (d DO) InstallNewSshKey(name, publicKey string) (string, error) {
	return "", errors.New("Not yet implemented!")
}

func (d DO) UpdateSshKey(id, name, key string) error {
	return errors.New("Not yet implemented!")
}

func (d *DO) GetConfig() *config.Config {
	return d.Config
}

func (d *DO) doGet(url string) (*http.Response, error) {
	cfg := d.GetConfig()
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", cfg.Providers[cfg.Provider].ApiKey))

	return client.Do(req)
}
