package config

import (
	"errors"
	"os"
	"strconv"
	"strings"

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
	Logging struct {
		TraceRequestsEnable bool `yaml:"traceRequestsEnable"`
	} `yaml:"logging"`
}

func Parse() (*Configuration, error) {
	configWay := os.Getenv("CONFIG")
	switch configWay {
	case "env":
		return parseFromENV()
	default:
		return parseFromFile()
	}
}

func parseFromFile() (*Configuration, error) {
	f, openFileErr := os.Open("configs/config.yml")

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

func parseFromENV() (*Configuration, error) {
	var validationErrors []string

	port := getEnv("SERVER_PORT", validationErrors)
	dbConnectionString := getEnv("DB_CONNECTION", validationErrors)
	dbName := getEnv("DB_NAME", validationErrors)
	jwtKey := getEnv("JWT_KEY", validationErrors)
	jwtIssuer := getEnv("JWT_ISSUER", validationErrors)
	jwtAudience := getEnv("JWT_AUDIENCE", validationErrors)
	jwtLifetimeDays := getEnv("JWT_LIFETIME_DAYS", validationErrors)
	hashKey := getEnv("HASH_KEY", validationErrors)

	if len(validationErrors) > 0 {
		return nil, errors.New(strings.Join(validationErrors, ","))
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}

	jwtLifetimeDaysInt, err := strconv.Atoi(jwtLifetimeDays)
	if err != nil {
		return nil, err
	}

	loggingTraceRequests := false
	traceRequests := os.Getenv("LOGGING_TRACE_REQUESTS")
	if traceRequests != "" {
		loggingTraceRequests, err = strconv.ParseBool(traceRequests)
		if err != nil {
			return nil, err
		}
	}

	conf := &Configuration{}
	conf.Server = struct {
		Port int `yaml:"port"`
	}(struct{ Port int }{Port: portInt})

	conf.Database = struct {
		ConnectionString string `yaml:"connectionString"`
		DbName           string `yaml:"dbName"`
	}(struct {
		ConnectionString string
		DbName           string
	}{ConnectionString: dbConnectionString, DbName: dbName})

	conf.JWT = struct {
		Key          string `yaml:"key"`
		Issuer       string `yaml:"issuer"`
		Audience     string `yaml:"audience"`
		LifetimeDays int    `yaml:"lifetimeDays"`
	}(struct {
		Key          string
		Issuer       string
		Audience     string
		LifetimeDays int
	}{Key: jwtKey, Issuer: jwtIssuer, Audience: jwtAudience, LifetimeDays: jwtLifetimeDaysInt})

	conf.Hash = struct {
		Key string `yaml:"key"`
	}(struct{ Key string }{Key: hashKey})

	conf.Logging = struct {
		TraceRequestsEnable bool `yaml:"traceRequestsEnable"`
	}(struct{ TraceRequestsEnable bool }{TraceRequestsEnable: loggingTraceRequests})

	return conf, nil
}

func getEnv(envName string, validationErrors []string) string {
	envVar := os.Getenv(envName)
	if envVar == "" {
		validationErrors = append(validationErrors, "env variable "+envVar+" is not set")
	}
	return envVar
}
