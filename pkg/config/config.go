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
	// MongoConnString string         `yaml:"mongo_conn_string"`
	MongoDB MongoDBConfig `yaml:"mongodb"`
}

type DatabaseConfig struct {
	File string `yaml:"file"`
}
type MongoDBConfig struct {
	ConnString string `yaml:"conn_string"`
	Name       string `yaml:"name"`
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
