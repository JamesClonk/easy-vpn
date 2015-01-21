package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/JamesClonk/easy-vpn/provider"
	"golang.org/x/crypto/ssh"
)

func sshExecCmd(p provider.API, ip string, cmd string) {
	key := readKeyFile(p.GetConfig().PrivateKeyFile)

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Println("Could not parse private key")
		log.Fatal(err)
	}

	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	client, err := ssh.Dial("tcp", ip+":22", config)
	if err != nil {
		log.Println("Could not connect to: " + ip)
		log.Fatal(err)
	}

	session, err := client.NewSession()
	if err != nil {
		log.Println("Could not create SSH session")
		log.Fatal(err)
	}
	defer session.Close()

	var stdOut bytes.Buffer
	var stdErr bytes.Buffer
	session.Stdout = &stdOut
	session.Stderr = &stdErr
	if err := session.Run(cmd); err != nil {
		log.Println("Could not run remote cmd through SSH: " + cmd)
		log.Println(stdErr.String())
		log.Fatal(err)
	}
	fmt.Println(stdOut.String())
}
