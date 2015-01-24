package ssh

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/JamesClonk/easy-vpn/provider"
	gossh "golang.org/x/crypto/ssh"
)

func Call(p provider.API, ip string, cmd string) {
	fmt.Println(Exec(p, ip, cmd))
}

func Exec(p provider.API, ip string, cmd string) string {
	out, err := Run(p, ip, cmd)
	if err != nil {
		log.Println("Could not run command through SSH: " + cmd)
		log.Fatal(err)
	}
	return out
}

func Run(p provider.API, ip string, cmd string) (string, error) {
	session := sshConnect(p, ip)
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

func WriteSelfDestruct(p provider.API, ip string, filename string) {
	session := sshConnect(p, ip)
	defer session.Close()

	data, err := ioutil.ReadFile(sanitizeFilename(filename))
	if err != nil {
		log.Println("Could not read in file: " + filename)
		log.Fatal(err)
	}

	go func() {
		writer, err := session.StdinPipe()
		if err != nil {
			log.Fatal(err)
		}
		defer writer.Close()

		fmt.Fprintln(writer, "C0750", len(data), "self-destruct.sh")
		fmt.Fprint(writer, data)
		fmt.Fprint(writer, "\x00")
	}()

	if err := session.Run("scp -qrt ./"); err != nil {
		log.Println("Could not transfer file through scp")
		log.Fatal(err)
	}
}

func sshConnect(p provider.API, ip string) *gossh.Session {
	key := readKeyFile(p.GetConfig().PrivateKeyFile)

	signer, err := gossh.ParsePrivateKey(key)
	if err != nil {
		log.Println("Could not parse private key")
		log.Fatal(err)
	}

	config := &gossh.ClientConfig{
		User: "root",
		Auth: []gossh.AuthMethod{
			gossh.PublicKeys(signer),
		},
	}

	client, err := gossh.Dial("tcp", ip+":22", config)
	if err != nil {
		log.Println("Could not connect to: " + ip)
		log.Fatal(err)
	}

	session, err := client.NewSession()
	if err != nil {
		log.Println("Could not create SSH session")
		log.Fatal(err)
	}
	return session
}
