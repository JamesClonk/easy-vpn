package main

import (
	"os"

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
			Name:  "provider, p",
			Value: "digitalocean",
			Usage: "specify which cloud VPS provider to use",
		},
		cli.StringFlag{
			Name:  "connect, c",
			Value: "yes",
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
			StartVpn(c)
		},
	}, {
		Name:        "down",
		ShortName:   "d",
		Usage:       "shutdown and destroy a VPS",
		Description: ".....", // TODO: add description, tell that it requires 1 argument: the VPS name/id to destroy
		Action: func(c *cli.Context) {
			DestroyVpn(c)
		},
	}, {
		Name:        "show",
		ShortName:   "s",
		Usage:       "shows all current VPN-VPS and their status",
		Description: ".....", // TODO: add description, will list all current VPS and their status that match naming criteria
		Action: func(c *cli.Context) {
			ShowVpn(c)
		},
	}}

	app.Action = func(c *cli.Context) {
		cli.ShowAppHelp(c)
		os.Exit(1)
	}
	app.RunAndExitOnError()
}

func StartVpn(c *cli.Context) {
}

func DestroyVpn(c *cli.Context) {
}

func ShowVpn(c *cli.Context) {
}
