package main

import (
	"log"

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

func loadConfiguration(filename string) *Config {
	var config Config
	if _, err := toml.DecodeFile(filename, &config); err != nil {
		// TODO: better error message & client aborting
		log.Fatal(err)
	}

	return &config
}
