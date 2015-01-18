package digitalocean

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/JamesClonk/easy-vpn/config"
	"github.com/JamesClonk/easy-vpn/provider"
)

var baseUrl = `https://api.digitalocean.com/v2`

type SshKeys struct {
	Keys []SshKey `json:"ssh_keys"`
}

type SshKey struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Key         string `json:"public_key"`
	Fingerprint string `json:"fingerprint"`
}

type Droplets struct {
	Droplets []Droplet `json:"droplets"`
}

type Droplet struct {
	Id     int      `json:"id"`
	Name   string   `json:"name"`
	Status string   `json:"status"`
	Region Region   `json:"region"`
	OS     Image    `json:"image"`
	IP     Networks `json:"networks"`
}

type Region struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type Image struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Distro string `json:"distribution"`
	Slug   string `json:"slug"`
}

type Networks struct {
	V4 []NetworkV4 `json:"v4"`
}

type NetworkV4 struct {
	IP   string `json:"ip_address"`
	Type string `json:"type"`
}

type DO struct {
	Config *config.Config
}

func (d DO) GetProviderName() string {
	return "digitalocean"
}

func (d DO) GetConfig() *config.Config {
	return d.Config
}

func (d DO) GetInstalledSshKeys() (data []provider.SshKey, err error) {
	resp, err := d.doGet(baseUrl + `/account/keys`)
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
	resp, err := d.doPost(baseUrl+`/account/keys`, values)
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

func (d DO) GetAllVMs() (data []provider.VM, err error) {
	resp, err := d.doGet(baseUrl + `/droplets`)
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

	var droplets Droplets
	if err := json.Unmarshal(body, &droplets); err != nil {
		return nil, err
	}

	// convert digitalocean droplets into array of provider api vm's
	for _, droplet := range droplets.Droplets {
		key := provider.VM{
			Id:     fmt.Sprintf("%d", droplet.Id),
			Name:   droplet.Name,
			Status: droplet.Status,
			OS:     droplet.OS.Slug,
			IP:     droplet.IP.V4[0].IP,
			Region: droplet.Region.Slug,
		}
		data = append(data, key)
	}

	return data, nil
}

func (d DO) CreateVM(name, os, size, region, sshkey string) (string, error) {
	log.Fatal("Not yet implemented!")
	return "", nil
}

func (d DO) StartVM(id string) error {
	log.Fatal("Not yet implemented!")
	return nil
}

func (d DO) Sleep() {
	time.Sleep(time.Duration(d.GetConfig().Sleep) * time.Millisecond)
}

func (d DO) deleteSshKey(id string) error {
	resp, err := d.doDelete(baseUrl + `/account/keys/` + id)
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

func (d *DO) doGet(url string) (*http.Response, error) {
	cfg := d.GetConfig()
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", cfg.Providers[cfg.Provider].ApiKey))

	d.Sleep() // respect request rate limitation
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

	d.Sleep() // respect request rate limitation
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

	d.Sleep() // respect request rate limitation
	return client.Do(req)
}
