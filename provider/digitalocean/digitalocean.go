package digitalocean

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

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
	defer resp.Body.Close()

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

func (d DO) InstallNewSshKey(name, key string) (string, error) {
	values := fmt.Sprintf(`{"name": "%v", "public_key": "%v"}`, name, key)
	resp, err := d.doPost(`https://api.digitalocean.com/v2/account/keys`, values)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusCreated {
		return "", errors.New(string(body))
	}

	if !strings.Contains(string(body), `"ssh_key":`) {
		return "", errors.New(string(body))
	}

	result := struct {
		Key SshKey `json:"ssh_key"`
	}{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", result.Key.Id), nil
}

func (d DO) UpdateSshKey(id, name, key string) (string, error) {
	// digitalocean has no "update" command, only delete and create
	// first we must destroy/delete the existing key
	err := d.deleteSshKey(id)
	if err != nil {
		return "", err
	}

	// then we can install a new key
	return d.InstallNewSshKey(name, key)
}

func (d DO) deleteSshKey(id string) error {
	resp, err := d.doDelete(`https://api.digitalocean.com/v2/account/keys/` + id)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return errors.New(string(body))
	}

	return nil
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

func (d *DO) doPost(url, data string) (*http.Response, error) {
	cfg := d.GetConfig()
	client := &http.Client{}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", cfg.Providers[cfg.Provider].ApiKey))
	req.Header.Set("Content-Type", "application/json")

	return client.Do(req)
}

func (d *DO) doDelete(url string) (*http.Response, error) {
	cfg := d.GetConfig()
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", cfg.Providers[cfg.Provider].ApiKey))

	return client.Do(req)
}
