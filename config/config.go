package config

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type ProberConfig map[string]interface{}

type Config struct {
	ProberConfigs []ProberConfig `yaml:"probers"`
}

var gConfig *Config = nil

func DefaultConfigPath() string {
	return "config.yml"
}

func GlobalConfig() *Config {

	if gConfig != nil {
		return gConfig
	}

	c := Config{}

	configPath := DefaultConfigPath()
	if configPath == "" {
		return &c
	}

	dat, err := ioutil.ReadFile(configPath)

	if err != nil {
		log.Println("Failed to load config file, err:", err)
		return &c
	}

	err = yaml.Unmarshal(dat, &c)

	if err != nil {
		log.Println("Failed to unmarshal config, err", err)
		return &c
	}

	gConfig = &c
	return gConfig
}
