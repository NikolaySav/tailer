package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

func NewConfig(path string) *Env {
	var c Env

	d, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatal(err, "Failed to read config file")
	}

	err = yaml.Unmarshal(d, &c)

	if err != nil {
		log.Fatal(err, "Failed to parse config file")
	}

	return &c
}

type Env struct {
	AuthType     string    `yaml:"authType"`
	DefaultLines int       `yaml:"defaultLines"`
	Projects     []Project `yaml:"projects"`
	Servers      []Server  `yaml:"servers"`
}

type Project struct {
	Name          string `yaml:"name"`
	FilePath      string `yaml:"filePath"`
	Server        string `yaml:"server"`
	BastionServer string `yaml:"bastionServer"`
}

type Server struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}
