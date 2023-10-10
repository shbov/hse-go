package Structs

import (
	"github.com/Masterminds/semver"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Port    int            `yaml:"port"`
	Version semver.Version `yaml:"version"`
}

func NewConfig(port int, v string) (*Config, error) {
	version, err := semver.NewVersion(v)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &Config{Port: port, Version: *version}, nil
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

	// todo: не очень понял, как можно парсить тип semver.Version
	var cfg struct {
		Port    int    `yaml:"port"`
		Version string `yaml:"version"`
	}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return NewConfig(cfg.Port, cfg.Version)
}
