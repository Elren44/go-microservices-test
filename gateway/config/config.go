package config

import "github.com/ilyakaznacheev/cleanenv"

type authConfig struct {
	DSN string `json:"dsn" env:"DSN" yaml:"dsn"`
	ENV string `json:"env" env:"ENV" yaml:"env"`
}

func NewAuthConfig() *authConfig {
	var authCfg authConfig
	err := cleanenv.ReadConfig("config.yml", &authCfg)
	if err != nil {
		panic(err)
	}

	return &authCfg
}
