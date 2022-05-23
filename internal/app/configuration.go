package app

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		ConnectionString string `yaml:"connectionString"`
		DbName           string `yaml:"dbName"`
	} `yaml:"database"`
	JWT struct {
		Key          string `yaml:"key"`
		Issuer       string `yaml:"issuer"`
		Audience     string `yaml:"audience"`
		LifetimeDays int    `yaml:"lifetimeDays"`
	} `yaml:"jwt"`
	Hash struct {
		Key string `yaml:"key"`
	} `yaml:"hash"`
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
