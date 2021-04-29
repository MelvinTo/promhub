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
var configPath = "config.yml"

func SetConfig(file string) {
	configPath = file
}

func GlobalConfig() *Config {

	if gConfig != nil {
		return gConfig
	}

	c := Config{}

	dat, err := ioutil.ReadFile(configPath)

	if err != nil {
		log.Fatal("Failed to load config file, err:", err)
		return &c
	}

	err = yaml.Unmarshal(dat, &c)

	if err != nil {
		log.Fatal("Failed to unmarshal config, err", err)
		return &c
	}

	gConfig = &c
	return gConfig
}
