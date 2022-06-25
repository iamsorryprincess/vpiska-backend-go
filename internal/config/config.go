package config

import (
	"github.com/caarlos0/env/v6"
)

type Configuration struct {
	ServerPort           int    `env:"SERVER_PORT" envDefault:"5000"`
	DbConnection         string `env:"DB_CONNECTION" envDefault:"mongodb://localhost:27017"`
	DbName               string `env:"DB_NAME" envDefault:"vpiska"`
	JWTKey               string `env:"JWT_KEY" envDefault:"vpiska_secretkey!123"`
	JWTIssuer            string `env:"JWT_ISSUER" envDefault:"VpiskaServer"`
	JWTAudience          string `env:"JWT_AUDIENCE" envDefault:"VpiskaClient"`
	JWTLifeTimeDays      int    `env:"JWT_LIFETIME_DAYS" envDefault:"3"`
	HashKey              string `env:"HASH_KEY" envDefault:"fbac497e4b44564f831f78d539b81a0c"`
	LoggingTraceRequests bool   `env:"LOGGING_TRACE_REQUESTS" envDefault:"false"`
}

func Parse() (*Configuration, error) {
	configuration := &Configuration{}

	err := env.Parse(configuration)
	if err != nil {
		return nil, err
	}

	return configuration, err
}
