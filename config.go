package main

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Provider  string              `toml:"provider"`
	Providers map[string]Provider `toml:"providers"`
	Options   Options             `toml:"options"`
}

type Provider struct {
	ApiKey string `toml:"api_key"`
	Region int    `toml:"region"`
}

type Options struct {
	Idletime    int  `toml:"max_idletime"`
	Uptime      int  `toml:"max_uptime"`
	Autoconnect bool `toml:"vpn_autoconnect"`
}

func loadConfiguration(filename string) (config *Config, err error) {
	if _, err = toml.DecodeFile(filename, &config); err != nil {
		return nil, err
	}
	return
}
