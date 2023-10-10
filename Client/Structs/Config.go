package Structs

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Host  string `yaml:"host"`
	Port  int    `yaml:"port"`
	Https bool   `yaml:"https"`
}

func ParseConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Error(err)
		}
	}(f)

	if err != nil {
		return nil, err
	}
	var cfg Config

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
