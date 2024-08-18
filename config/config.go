package config

import (
	"github.com/naoina/toml"
	"os"
)

type Config struct {
	DB struct {
		Database string
		URL      string
	}

	Kafka struct {
		URL      string
		ClientID string
	}
}

func NewConfig(path string) *Config {
	c := new(Config)
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	if err = toml.NewDecoder(f).Decode(c); err != nil {
		panic(err)
	}
	return c
}
