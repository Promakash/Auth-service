package config

import (
	"auth_service/pkg/infra"
	pkglog "auth_service/pkg/log"
)

type HTTPConfig struct {
	Address string               `yaml:"address"`
	PG      infra.PostgresConfig `yaml:"postgres"`
	Log     pkglog.Config        `yaml:"log"`
	Jwt     JWTConfig            `yaml:"jwt"`
	SMTP    infra.SMTPConfig     `yaml:"smtp"`
}

type JWTConfig struct {
	Secret     string `yaml:"secret"`
	RefreshExp int    `yaml:"refreshExp"`
	AccessExp  int    `yaml:"accessExp"`
}
