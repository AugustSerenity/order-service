package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server Server `yaml:"server"`
	DB     DB     `yaml:"db"`
}

type Server struct {
	Address         string        `yaml:"address"`
	Timeout         time.Duration `yaml:"timeout"`
	IdleTimeout     time.Duration `yaml:"idle_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type DB struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port" env-default:"5432"`
	Username string `yaml:"username"`
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
}

func ParseConfig(path string) *Config {
	var cfg Config

	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file, %s", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	return &cfg
}
