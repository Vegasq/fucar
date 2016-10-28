package main

import (
	"bytes"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

type SSHNode struct {
	host     string
	port     string
	username string
	password string

	debug bool
}

func (node *SSHNode) Exec(command string) string {
	session := node.connect()

	var b bytes.Buffer
	session.Stdout = &b

	if err := session.Run(command); err != nil {
		log.Fatal("Failed to run: " + command + " : " + err.Error())
	}

	session.Close()
	return b.String()
}

func (node *SSHNode) getDial() *ssh.Client {
	config := &ssh.ClientConfig{
		User: node.username,
		Auth: []ssh.AuthMethod{
			ssh.Password(node.password),
		},
	}

	client, err := ssh.Dial("tcp", node.host+":"+node.port, config)
	if err != nil {
		if node.debug {
			log.Println("Failed to dial host: " + node.host + ". Sleep for 1s and retry.")
		}
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
