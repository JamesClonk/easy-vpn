package vultr

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/JamesClonk/easy-vpn/config"
	"github.com/JamesClonk/easy-vpn/provider"
)

type SshKey struct {
	Id          string    `json:"SSHKEYID"`
	Name        string    `json:"name"`
	Key         string    `json:"ssh_key"`
	CreatedDate time.Time `json:"date_created"`
}

type Vultr struct {
	Config *config.Config
}

func (v Vultr) GetInstalledSshKeys() (data []provider.SshKey, err error) {
	resp, err := http.Get(v.urlWithApiKey(`https://api.vultr.com/v1/sshkey/list`))
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

	// vultr returns empty array if no SSH Keys are found
	if string(body) == "[]" {
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

func (v *Vultr) GetConfig() *config.Config {
	return v.Config
}

func (v *Vultr) urlWithApiKey(url string) string {
	cfg := v.GetConfig()
	return fmt.Sprintf("%v?api_key=%v", url, cfg.Providers[cfg.Provider].ApiKey)
}
