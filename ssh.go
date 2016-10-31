package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type SSHNode struct {
	host     string
	username string
}

func (node *SSHNode) Exec(command string) string {
	session := node.connect()

	output, err := session.Output(command)
	output_str := string(output)

	if err != nil {
		if strings.Contains(output_str, "local node is not a member of the token ring") {
			log.Println("Node already removed")
		} else {
			log.Println(output_str)
			log.Fatal("Failed to run: " + command + " : " + err.Error())
			os.Exit(2)
		}
	}

	session.Close()
	return output_str
}

func (node *SSHNode) getDial() *ssh.Client {
	key, err := ioutil.ReadFile("/root/.ssh/id_rsa")
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User: node.username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	client, err := ssh.Dial("tcp", node.host+":22", config)
	if err != nil {
		log.Println("Failed to dial host: " + node.host + ". Sleep for 1s and retry.")
		log.Fatalln(err)
		time.Sleep(1000)
		return node.getDial()
	}
	return client
}

func (node *SSHNode) connect() *ssh.Session {
	client := node.getDial()
	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
		os.Exit(1)
	}
	return session
}
