package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/JamesClonk/easy-vpn/config"
	"github.com/JamesClonk/easy-vpn/provider"
	"github.com/JamesClonk/easy-vpn/provider/digitalocean"
	"github.com/JamesClonk/easy-vpn/provider/vultr"
	"github.com/codegangsta/cli"
)

const (
	VERSION            = "1.0.0"
	EASYVPN_IDENTIFIER = "easy-vpn"
)

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
			Name:  "api-key, k",
			Value: "abc123xyz",
			Usage: "API-Key for cloud VPS provider",
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
	p := getProvider(c)

	sshkeyId := getEasyVpnSshKeyId(p)
	vm := getEasyVpnVM(p, sshkeyId)

	fmt.Printf("%q\n", vm) // TODO: prettify

	// update machine
	sshExecCmd(p, vm.IP, `apt-get update -qq`)
	sshExecCmd(p, vm.IP, `apt-get install -qy docker pptpd iptables`)

	// setup docker
	sshExecCmd(p, vm.IP, `docker pull JamesClonk/docker-pptpd`)
	sshExecCmd(p, vm.IP, `echo "easy-vpn * secret-password *" >> /chap-secrets`)

	// run docker
	sshExecCmd(p, vm.IP, `docker run --name pptpd --privileged -d -p 1723:1723 -v /chap-secrets:/etc/ppp/chap-secrets:ro JamesClonk/docker-pptpd`)
}

func destroyVpn(c *cli.Context) {
	p := getProvider(c)
	destroyEasyVpnVM(p)
}

func showVpn(c *cli.Context) {
	p := getProvider(c)
	for _, vm := range getAllVMs(p) {
		fmt.Printf("%q\n", vm) // TODO: prettify
	}
}

func parseGlobalOptions(c *cli.Context) *config.Config {
	cfg, err := config.LoadConfiguration(c.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}

	if c.GlobalIsSet("provider") {
		cfg.Provider = c.GlobalString("provider")
	}

	if c.GlobalIsSet("api-key") {
		cfg.Providers[cfg.Provider] = config.Provider{
			ApiKey: c.GlobalString("api-key"),
			Region: cfg.Providers[cfg.Provider].Region,
			Size:   cfg.Providers[cfg.Provider].Size,
			OS:     cfg.Providers[cfg.Provider].OS,
		}
	}

	if c.GlobalIsSet("autoconnect") {
		cfg.Options.Autoconnect = strings.ToLower(c.GlobalString("autoconnect")) == "true"
	}

	if c.GlobalIsSet("idletime") {
		idletime, err := strconv.ParseInt(c.GlobalString("idletime"), 10, 32)
		if err != nil {
			log.Fatalf("Invalid value for --idletime option given: %v\n", c.GlobalString("idletime"))
		}
		cfg.Options.Idletime = int(idletime)
	}

	if c.GlobalIsSet("uptime") {
		uptime, err := strconv.ParseInt(c.GlobalString("uptime"), 10, 32)
		if err != nil {
			log.Fatalf("Invalid value for --uptime option given: %v\n", c.GlobalString("uptime"))
		}
		cfg.Options.Uptime = int(uptime)
	}

	return cfg
}

func getProvider(c *cli.Context) provider.API {
	cfg := parseGlobalOptions(c)

	switch cfg.Provider {
	case "digitalocean":
		return digitalocean.DO{Config: cfg}
	case "vultr":
		return vultr.Vultr{Config: cfg}
	case "aws":
		log.Fatal("Not yet implemented!")
	default:
		log.Fatal("Unknown provider!")
	}
	return nil
}
