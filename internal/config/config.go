package config

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

func Parse() (*Configuration, error) {
	env := os.Getenv("APP_ENV")
	configPath := ""

	switch env {
	case "prod":
		configPath = "configs/prod.yml"
		break
	default:
		configPath = "configs/test.yml"
		break
	}

	f, openFileErr := os.Open(configPath)

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
