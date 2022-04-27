package app

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		ConnectionString string `yaml:"connectionString"`
		DbName           string `yaml:"dbName"`
	} `yaml:"database"`
}

func parseConfig() (*Configuration, error) {
	f, openFileErr := os.Open("configs/main.yml")

	if openFileErr != nil {
		return nil, openFileErr
	}
	defer f.Close()

	conf := &Configuration{}
	decoder := yaml.NewDecoder(f)
	decodeErr := decoder.Decode(conf)

	if decodeErr != nil {
		return nil, decodeErr
	}

	return conf, nil
}
