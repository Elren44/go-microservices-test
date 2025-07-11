package config

import "github.com/ilyakaznacheev/cleanenv"

type AuthConfig struct {
	DSN    string `json:"dsn" env:"dsn" yaml:"dsn"`
	ENV    string `json:"env" env:"env" yaml:"env"`
	Secret string `json:"secret" env:"secret" yaml:"secret"`
}

func NewAuthConfig() *AuthConfig {
	var authCfg AuthConfig
	err := cleanenv.ReadConfig("config.yml", &authCfg)
	if err != nil {
		panic(err)
	}

	return &authCfg
}
