package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"

	"github.com/JamesClonk/easy-vpn/provider"
	"golang.org/x/crypto/ssh"
)

func sshCall(p provider.API, ip string, cmd string) {
	fmt.Println(sshExec(p, ip, cmd))
}

func sshExec(p provider.API, ip string, cmd string) string {
	out, err := sshRun(p, ip, cmd)
	if err != nil {
		log.Println("Could not run command through SSH: " + cmd)
		log.Fatal(err)
	}
	return out
}

func sshRun(p provider.API, ip string, cmd string) (string, error) {
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
		return "", errors.New(fmt.Sprintf("%s\n%v", stdErr.String(), err))
	}

	return stdOut.String(), nil
}
