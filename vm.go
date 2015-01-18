package main

import (
	"log"

	"github.com/JamesClonk/easy-vpn/provider"
)

func getAllVMs(p provider.API) []provider.VM {
	machines, err := p.GetAllVMs()
	if err != nil {
		log.Println("Could not retrieve list of virtual machines")
		log.Fatal(err)
	}
	return machines
}
