package main

import (
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"os"
)

func main() {

	login := "test"
	password := "SDHBCXdsedfs222"
	ip := "151.248.113.144"
	port := "443"

	config := &ssh.ClientConfig{
		User: login,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", ip+":"+port, config)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	sess, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer sess.Close()

	stdin, err := sess.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	go io.Copy(stdin, os.Stdin)

	sess.Stdin = os.Stdin
	sess.Stdout = os.Stdout
	sess.Stderr = os.Stderr

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := sess.RequestPty("xterm", 40, 80, modes); err != nil {
		log.Fatal(err)
	}

	if err := sess.Shell(); err != nil {
		log.Fatal(err)
	}

	if err = sess.Wait(); err != nil {
		log.Fatal(err)
	}

}
