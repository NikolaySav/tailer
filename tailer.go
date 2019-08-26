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

	servers := make(map[string]config.Server)
	projects := make(map[string]config.Project)

	fmt.Println(c)

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

	fmt.Println("Project: ", project)

	server, exists := servers[project.Server]

	if !exists {
		log.Fatal("Project config not found")
	}

	fmt.Println("Server: ", server.Name)

	sshConfig := &ssh.ClientConfig{
		User: server.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(server.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", "192.168.2.114:22", sshConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	conn := &remote.Connection{Client: client}

	command := fmt.Sprintf("tail -f -n%d %s", lines, project.FilePath)
	fmt.Println("Executing: ", command)

	conn.Run(command)
}
