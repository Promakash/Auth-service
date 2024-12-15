package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config interface{}

func ParseAppConfig[T Config](env string) *T {
	configPath := os.Getenv(env)
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}

	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("error opening config file: %v", err)
	}

	cfg, err := Bind[T](configPath)
	if err != nil {
		log.Fatalf("can't unmarshal application config: %v", err)
	}

	return &cfg
}

func Bind[T Config](configPath string) (T, error) {
	var cfg T

	err := cleanenv.ReadConfig(configPath, &cfg)
	return cfg, err
}
