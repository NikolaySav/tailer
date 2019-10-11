package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"tailer/internal/config"
	"tailer/internal/remote"

	"golang.org/x/crypto/ssh"
)

var servers = make(map[string]config.Server)

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		log.Fatal("Provide project name as an argument")
	}

	projectName := os.Args[1]
	c := config.NewConfig("./config.yaml")

	var lines int
	var err error

	if len(os.Args) == 3 {
		lines, err = strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal("Invalid lines number value")
		}
	} else {
		lines = c.DefaultLines
	}

	projects := make(map[string]config.Project)

	for _, project := range c.Projects {
		projects[project.Name] = project
	}

	for _, server := range c.Servers {
		servers[server.Name] = server
	}

	project, exists := projects[projectName]
	if !exists {
		log.Fatal("Project config not found")
	}

	var bastionClient, client *ssh.Client

	// resolve ssh client
	if project.BastionServer != "" {
		bastionClient, client = newBastionHostSshClient(project)
		defer bastionClient.Close()
	} else {
		client = newSshClient(project)
	}
	defer client.Close()

	targetConn := remote.Connection{Client: client}

	command := fmt.Sprintf("tail -f -n%d %s", lines, project.FilePath)
	fmt.Println("Executing: ", command)

	targetConn.Run(command)
}

func newBastionHostSshClient(project config.Project) (*ssh.Client, *ssh.Client) {
	bastionServer, exists := servers[project.BastionServer]

	if !exists {
		log.Fatal("BastionServer config not found")
	}

	fmt.Println("Bastion Server: ", bastionServer.Name)

	targetServer, exists := servers[project.Server]

	if !exists {
		log.Fatal("Target server config not found")
	}

	fmt.Println("Target Server: ", targetServer.Name)

	sshConfig := &ssh.ClientConfig{
		User: bastionServer.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(bastionServer.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// ssh to bastion server
	bastionClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", bastionServer.Host, bastionServer.Port), sshConfig)
	if err != nil {
		log.Fatal(err)
	}

	targetServerAddr := fmt.Sprintf("%s:%s", targetServer.Host, targetServer.Port)

	// connection to target server
	conn, err := bastionClient.Dial("tcp", targetServerAddr)
	if err != nil {
		log.Fatal(err)
	}

	ncc, chans, reqs, err := ssh.NewClientConn(conn, targetServerAddr, &ssh.ClientConfig{
		User: targetServer.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(targetServer.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		log.Fatal(err)
	}

	targetClient := ssh.NewClient(ncc, chans, reqs)

	return bastionClient, targetClient
}

func newSshClient(project config.Project) *ssh.Client {
	targetServer, exists := servers[project.Server]

	if !exists {
		log.Fatal("Target server config not found")
	}

	fmt.Println("Target Server: ", targetServer.Name)

	sshConfig := &ssh.ClientConfig{
		User: targetServer.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(targetServer.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", targetServer.Host, targetServer.Port), sshConfig)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
