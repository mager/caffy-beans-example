package config

import (
	"io/ioutil"

	"go.uber.org/fx"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Application `yaml:"application"`
}

type Application struct {
	Address string `yaml:"address"`
}

func ProvideConfig() *Config {
	conf := Config{}
	data, err := ioutil.ReadFile("config/base.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal([]byte(data), &conf)
	if err != nil {
		panic(err)
	}

	return &conf
}

var Module = fx.Provide(ProvideConfig)
