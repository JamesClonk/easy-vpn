package vultr

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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

type Server struct {
	Id     string `json:"SUBID"`
	Name   string `json:"label"`
	OS     string `json:"os"`
	IP     string `json:"main_ip"`
	Region string `json:"DCID"`
	Status string `json:"status"`
}

type Vultr struct {
	Config *config.Config
}

func (v Vultr) GetProviderName() string {
	return "VULTR"
}

func (v Vultr) GetConfig() *config.Config {
	return v.Config
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

func (v Vultr) GetAllVMs() (data []provider.VM, err error) {
	resp, err := http.Get(v.urlWithApiKey(baseUrl + `/server/list`))
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

	// vultr returns empty array if no servers are found
	if strings.Trim(string(body), "\t\r\n ") == "[]" {
		return data, nil
	}

	var vultrServers map[string]Server
	if err := json.Unmarshal(body, &vultrServers); err != nil {
		return nil, err
	}

	// convert vultr servers into array of provider api vm's
	for _, value := range vultrServers {
		key := provider.VM{
			Id:     value.Id,
			Name:   value.Name,
			OS:     value.OS,
			IP:     value.IP,
			Region: value.Region,
			Status: value.Status,
		}
		data = append(data, key)
	}

	return data, nil
}

func (v Vultr) CreateVM(name, os, size, region string) (string, error) {
	log.Fatal("Not yet implemented!")
	return "", nil
}

func (v *Vultr) urlWithApiKey(url string) string {
	cfg := v.GetConfig()
	return fmt.Sprintf("%v?api_key=%v", url, cfg.Providers[cfg.Provider].ApiKey)
}
