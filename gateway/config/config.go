package config

import "github.com/ilyakaznacheev/cleanenv"

type GatewayConfig struct {
	DSN            string `json:"dsn" env:"DSN" yaml:"dsn"`
	ENV            string `json:"env" env:"ENV" yaml:"env"`
	AuthServiceURL string `json:"auth_service_url" env:"AUTH_SERVICE_URL" yaml:"auth_service_url"`
	Secret         string `json:"secret" env:"SECRET" yaml:"secret"`
}

func NewGatewayConfig() *GatewayConfig {
	var gatewayCfg GatewayConfig
	err := cleanenv.ReadConfig("config.yml", &gatewayCfg)
	if err != nil {
		panic(err)
	}

	return &gatewayCfg
}
