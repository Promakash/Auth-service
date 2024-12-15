package config

import (
	"auth_service/pkg/infra"
	pkglog "auth_service/pkg/log"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type HTTPConfig struct {
	Address    string               `yaml:"address"`
	PG         infra.PostgresConfig `yaml:"postgres"`
	Log        pkglog.Config        `yaml:"log"`
	Secret     string               `yaml:"secret"`
	RefreshExp int                  `yaml:"refreshExp"`
	AccessExp  int                  `yaml:"accessExp"`
}

func MustLoad() *HTTPConfig {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}

	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("error opening config file: %s", err)
	}

	var cfg HTTPConfig

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	return &cfg
}
