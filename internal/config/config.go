package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Host           string
	Postgresqlport int
	User           string
	Password       string
	DBname         string
	Port           string
	ConnectionType string
}

func GetConfig() (*Config, error) {
	var err error
	var conf *Config

	_, err = toml.DecodeFile("./configs/config.toml", &conf)

	return conf, err
}
