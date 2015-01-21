package main

import (
	"io/ioutil"
	"log"
	"os/user"
	"strings"

	"github.com/JamesClonk/easy-vpn/provider"
)

func getEasyVpnSshKeyId(p provider.API) (keyId string) {
	key := string(readKeyFile(p.GetConfig().PublicKeyFile))

	// first lets get all currently installed ssh-keys
	keys, err := p.GetInstalledSshKeys()
	if err != nil {
		log.Println("Could not retrieve list of installed SSH-Keys")
		log.Fatal(err)
	}

	// then check to see if easy-vpn ssh-key is already installed
	keyInstalled := false
	for _, key := range keys {
		if key.Name == EASYVPN_IDENTIFIER {
			keyId = key.Id
			keyInstalled = true
			break
		}
	}

	// if it is already installed, update it to make sure its public-key is up-to-date
	if keyInstalled {
		keyId, err = p.UpdateSshKey(keyId, EASYVPN_IDENTIFIER, key)
		if err != nil {
			log.Println("Could not update SSH-Key")
			log.Fatal(err)
		}
	} else { // otherwise, install as a new ssh-key
		keyId, err = p.InstallNewSshKey(EASYVPN_IDENTIFIER, key)
		if err != nil {
			log.Println("Could not install SSH-Key")
			log.Fatal(err)
		}
	}

	return keyId
}

func readKeyFile(filename string) []byte {
	data, err := ioutil.ReadFile(sanitizeFilename(filename))
	if err != nil {
		log.Println("Could not read ssh key file: " + filename)
		log.Fatal(err)
	}
	return data
}

func sanitizeFilename(filename string) string {
	// replace beginning tilde (~) character with path to users home directory
	if strings.HasPrefix(filename, `~`) {
		usr, err := user.Current()
		if err != nil {
			log.Println("Could not get information about current user")
			log.Fatal(err)
		}
		home := usr.HomeDir
		filename = strings.Replace(filename, `~`, home, 1)
	}

	return filename
}
