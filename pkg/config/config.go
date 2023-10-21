package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

const configFile = "config.yaml"

var Config *Configuration

type Configuration struct {
	Environment string         `yaml:"environment"`
	Database    DatabaseConfig `yaml:"database"`
	SigningKey  string         `yaml:"signing_key"`
}

type DatabaseConfig struct {
	File string `yaml:"file"`
}

func loadConfigFile() {

	file, err := os.Open(configFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&Config)
	if err != nil {
		panic(err)
	}
}

func init() {
	loadConfigFile()
}
