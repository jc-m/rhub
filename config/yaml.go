package config

import (
	"github.com/go-yaml/yaml"
	"log"
	"io/ioutil"
)

type Module struct {
	Name string
	Module string
	Upstream string
	Config map[string]string
}

type Stream struct {
	Name string
	Modules []Module
}

type Config struct {
	Streams []Stream
}

func GetConfig(file string) *Config {
	var data []byte

	v := &Config{}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Cannot read file: %v", err)
	}

	err = yaml.Unmarshal(data, v)
	if err != nil {
		log.Fatalf("Error Parsing config: %v", err)
	}

	return v
}