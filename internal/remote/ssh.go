package remote

import (
	"io"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

type Connection struct {
	*ssh.Client
}

func (c Connection) Run(cmd string) {
	sess, err := c.Client.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()

	sessStdOut, err := sess.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	go io.Copy(os.Stdout, sessStdOut)

	// todo лишнее?
	sessStderr, err := sess.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	go io.Copy(os.Stderr, sessStderr)

	err = sess.Run(cmd)
	if err != nil {
		log.Fatal(err)
	}
}
