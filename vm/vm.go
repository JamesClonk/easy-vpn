package vm

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/JamesClonk/easy-vpn/provider"
	"github.com/JamesClonk/easy-vpn/ssh"
)

func GetAll(p provider.API) []provider.VM {
	machines, err := p.GetAllVMs()
	if err != nil {
		log.Println("Could not retrieve list of virtual machines")
		log.Fatal(err)
	}
	return machines
}

func GetEasyVpn(p provider.API, sshkeyId string, vmName string) (vm provider.VM) {
	cfg := p.GetConfig()
	os := cfg.Providers[cfg.Provider].OS
	size := cfg.Providers[cfg.Provider].Size
	region := cfg.Providers[cfg.Provider].Region

	// check to see if easy-vpn vm already exists
	vmExists := false
	for _, machine := range GetAll(p) {
		if machine.Name == vmName {
			vm = machine
			vmExists = true
			break
		}
	}

	if vmExists {
		fmt.Println("Virtual machine already exists")
	} else { // create a new vm and start it if it did not yet exist
		fmt.Println("Create new virtual machine")

		_, err := p.CreateVM(vmName, os, size, region, sshkeyId)
		if err != nil {
			log.Println("Could not create new virtual machine")
			log.Fatal(err)
		}
		waitForNewVM(p, &vm, vmName)
	}

	// make sure its up and running
	statusOfVM(p, &vm)
	// this is needed because some providers such as vultr do a dist-upgrade on new vms
	readynessOfVM(p, &vm)

	fmt.Println()

	return
}

func DestroyEasyVpn(p provider.API, vmName string) {
	var vm provider.VM

	// check to see if easy-vpn vm actually exists
	vmExists := false
	for _, machine := range GetAll(p) {
		if machine.Name == vmName {
			vm = machine
			vmExists = true
			break
		}
	}

	// ask to destroy it if it exists
	if vmExists {
		fmt.Println("Do you really want to destroy the following virtual machine?")
		fmt.Printf("%q\n", vm)
		fmt.Printf(`Confirm with "YES": `)

		reader := bufio.NewReader(os.Stdin)
		answer, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		answer = strings.Trim(answer, "\t\n\r ")

		if answer == "YES" {
			fmt.Println("Destroy virtual machine")
			err := p.DestroyVM(vm.Id)
			if err != nil {
				log.Println("Could not destroy virtual machine")
				log.Fatal(err)
			}
		}
	} else {
		fmt.Println("Virtual machine did not exist")
	}

}

func waitForNewVM(p provider.API, vm *provider.VM, vmName string) {
	fmt.Printf("Virtual machine installation")
	ticker := ticker()

	// TODO: maybe have some maximum waiting/polling time for doing a timeout
POLL:
	for {
		for _, machine := range GetAll(p) {
			if machine.Name == vmName {
				vm.Id = machine.Id
				vm.Name = machine.Name
				vm.Status = machine.Status
				vm.IP = machine.IP
				vm.OS = machine.OS
				vm.Region = machine.Region

				ticker.Stop()
				break POLL
			}
		}
		time.Sleep(15 * time.Second)
	}
	fmt.Printf("\nVirtual machine created: %q\n", vm) // TODO: prettify
}

func statusOfVM(p provider.API, vm *provider.VM) {
	fmt.Printf("Virtual machine status check")
	ticker := ticker()

	// TODO: maybe have some maximum waiting/polling time for doing a timeout
POLL:
	for {
		for _, machine := range GetAll(p) {
			if machine.Id == vm.Id &&
				machine.Status == "active" {
				vm.Status = machine.Status
				vm.IP = machine.IP

				ticker.Stop()
				break POLL
			}
		}
		time.Sleep(10 * time.Second)
	}
	fmt.Printf("\nVirtual machine is active\n")
}

func readynessOfVM(p provider.API, vm *provider.VM) {
	fmt.Printf("Virtual machine readyness check")
	ticker := ticker()

	// TODO: maybe have some maximum waiting/polling time for doing a timeout
POLL:
	for {
		// TODO: improve apt-get lock check
		out := ssh.Exec(p, vm.IP, `lsof /var/lib/dpkg/lock >/dev/null 2>&1; [ $? = 0 ] && echo "locked"; echo "..."`)
		if !strings.Contains(out, "locked") {
			ticker.Stop()
			break POLL
		}
		time.Sleep(15 * time.Second)
	}
	fmt.Printf("\nVirtual machine is ready\n")
}

func ticker() *time.Ticker {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			fmt.Print(".")
		}
	}()
	return ticker
}
