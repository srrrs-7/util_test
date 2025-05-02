package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	FILE_PATH = "/tmp/proxy/config.yaml"
)

type Config struct {
	Server struct {
		User string `yaml:"user"`
		Pass string `yaml:"pass"`
		Addr string `yaml:"addr"`
	} `yaml:"server"`
	Database struct {
		Test struct {
			User string `yaml:"user"`
			Pass string `yaml:"pass"`
			Addr string `yaml:"addr"`
			Name string `yaml:"name"`
		}
	} `yaml:"database"`
	Setting struct {
		MinOpenConns int `yaml:"min_open_conns"`
		MaxIdleConns int `yaml:"max_idle_conns"`
		MaxOpenConns int `yaml:"max_open_conns"`
		ConnLifetime int `yaml:"conn_lifetime"`
	}
}

func NewConfig() Config {
	data, err := os.ReadFile(FILE_PATH)
	if err != nil {
		log.Panicf("Failed to read config file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Panicf("Failed to unmarshal config file: %v", err)
	}

	if config.Server.User == "" ||
		config.Server.Pass == "" ||
		config.Server.Addr == "" ||
		config.Database.Test.User == "" ||
		config.Database.Test.Pass == "" ||
		config.Database.Test.Addr == "" ||
		config.Database.Test.Name == "" ||
		config.Setting.MinOpenConns <= 0 ||
		config.Setting.MaxIdleConns <= 0 ||
		config.Setting.MaxOpenConns <= 0 ||
		config.Setting.ConnLifetime <= 0 {
		log.Panicf("Server configuration is missing user, pass, or addr")
	}

	return config
}
