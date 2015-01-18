package vultr

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/JamesClonk/easy-vpn/config"
	"github.com/JamesClonk/easy-vpn/provider"
)

var baseUrl = `https://api.vultr.com/v1`

type SshKey struct {
	Id   string `json:"SSHKEYID"`
	Name string `json:"name"`
	Key  string `json:"ssh_key"`
}

type Vultr struct {
	Config *config.Config
}

func (v Vultr) GetProviderName() string {
	return "VULTR"
}

func (v Vultr) GetInstalledSshKeys() (data []provider.SshKey, err error) {
	resp, err := http.Get(v.urlWithApiKey(baseUrl + `/sshkey/list`))
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

	// vultr returns empty array if no SSH Keys are found
	if strings.Trim(string(body), "\t\r\n ") == "[]" {
		return data, nil
	}

	var vultrKeys map[string]SshKey
	if err := json.Unmarshal(body, &vultrKeys); err != nil {
		return nil, err
	}

	// convert vultr ssh-keys into array of provider api ssh-keys
	for _, value := range vultrKeys {
		key := provider.SshKey{
			Id:   value.Id,
			Name: value.Name,
			Key:  value.Key,
		}
		data = append(data, key)
	}

	return data, nil
}

func (v Vultr) InstallNewSshKey(name, key string) (string, error) {
	resp, err := http.PostForm(v.urlWithApiKey(baseUrl+`/sshkey/create`),
		url.Values{
			"name":    {name},
			"ssh_key": {key},
		})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(string(body))
	}

	if !strings.Contains(string(body), `"SSHKEYID":`) {
		return "", errors.New(string(body))
	}

	result := struct {
		Id string `json:"SSHKEYID"`
	}{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return result.Id, nil
}

func (v Vultr) UpdateSshKey(id, name, key string) (string, error) {
	resp, err := http.PostForm(v.urlWithApiKey(baseUrl+`/sshkey/update`),
		url.Values{
			"SSHKEYID": {id},
			"name":     {name},
			"ssh_key":  {key},
		})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(string(body))
	}

	return id, nil
}

func (v Vultr) GetConfig() *config.Config {
	return v.Config
}

func (v *Vultr) urlWithApiKey(url string) string {
	cfg := v.GetConfig()
	return fmt.Sprintf("%v?api_key=%v", url, cfg.Providers[cfg.Provider].ApiKey)
}
