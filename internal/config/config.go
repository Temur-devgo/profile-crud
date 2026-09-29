package config

import (
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	"os"
)

var conf Config

const configFileName = "config.yaml"

type Config struct {
	ServerHost string `yaml:"server_host"`
	ServerPort string `yaml:"server_port"`
	DB         string `yaml:"db"`
}

func GetConfigs() Config {
	yamlFile, err := os.ReadFile(configFileName)
	if err != nil {
		log.Fatal().Msgf("config file '%s' not exist", configFileName)
	}

	if err = yaml.Unmarshal(yamlFile, &conf); err != nil {
		log.Fatal().Msg("config file is not correct format")
	}

	return conf
}
