package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/codegangsta/cli"
)

const VERSION = "1.0.0"

func main() {
	app := cli.NewApp()
	app.Name = "easy-vpn"
	app.Author = "JamesClonk"
	app.Email = "jamesclonk@jamesclonk.ch"
	app.Version = VERSION
	app.Usage = "a simple tool to spin up a VPN server on a cloud VPS that self-destructs after idle time"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "easy-vpn.toml",
			Usage: "specify which configuration file to use",
		},
		cli.StringFlag{
			Name:  "provider, p",
			Value: "digitalocean",
			Usage: "specify which cloud VPS provider to use",
		},
		cli.StringFlag{
			Name:  "autoconnect, a",
			Value: "true",
			Usage: "do automatic VPN connect after a VPS was started?",
		},
		cli.StringFlag{
			Name:  "idletime, i",
			Value: "15",
			Usage: "idle time in minutes after which the VPS will self-destruct",
		},
		cli.StringFlag{
			Name:  "uptime, u",
			Value: "360",
			Usage: "maximum uptime in minutes after which the VPS will self-destruct",
		},
	}

	app.Commands = []cli.Command{{
		Name:        "up",
		ShortName:   "u",
		Usage:       "spin up a new VPS with a VPN server in it",
		Description: ".....", // TODO: add description, explain -r/--region option
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "region, r",
				Usage: "specify which region to use for new VPS",
			},
		},
		Action: func(c *cli.Context) {
			startVpn(c)
		},
	}, {
		Name:        "down",
		ShortName:   "d",
		Usage:       "shutdown and destroy a VPS",
		Description: ".....", // TODO: add description, tell that it requires 1 argument: the VPS name/id to destroy
		Action: func(c *cli.Context) {
			destroyVpn(c)
		},
	}, {
		Name:        "show",
		ShortName:   "s",
		Usage:       "shows all current VPN-VPS and their status",
		Description: ".....", // TODO: add description, will list all current VPS and their status that match naming criteria
		Action: func(c *cli.Context) {
			showVpn(c)
		},
	}}

	app.Action = func(c *cli.Context) {
		cli.ShowAppHelp(c)
		os.Exit(1)
	}
	app.RunAndExitOnError()
}

func startVpn(c *cli.Context) {
	//config := parseGlobalOptions(c)
}

func destroyVpn(c *cli.Context) {
	//config := parseGlobalOptions(c)
}

func showVpn(c *cli.Context) {
	//config := parseGlobalOptions(c)
}

func parseGlobalOptions(c *cli.Context) *Config {
	config, err := loadConfiguration(c.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}

	if c.GlobalIsSet("provider") {
		config.Provider = c.GlobalString("provider")
	}

	if c.GlobalIsSet("autoconnect") {
		config.Options.Autoconnect = strings.ToLower(c.GlobalString("autoconnect")) == "true"
	}

	if c.GlobalIsSet("idletime") {
		idletime, err := strconv.ParseInt(c.GlobalString("idletime"), 10, 32)
		if err != nil {
			log.Fatalf("Invalid value for --idletime option given: %v\n", c.GlobalString("idletime"))
		}
		config.Options.Idletime = int(idletime)
	}

	if c.GlobalIsSet("uptime") {
		uptime, err := strconv.ParseInt(c.GlobalString("uptime"), 10, 32)
		if err != nil {
			log.Fatalf("Invalid value for --uptime option given: %v\n", c.GlobalString("uptime"))
		}
		config.Options.Uptime = int(uptime)
	}

	return config
}
