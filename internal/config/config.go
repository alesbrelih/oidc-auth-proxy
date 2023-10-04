package config

import (
	"github.com/caarlos0/env/v9"
)

type Config struct {
	ClientID           string `env:"GOAP_CLIENT_ID"`
	ClientSecret       string `env:"GOAP_CLIENT_SECRET"`
	Issuer             string `env:"GOAP_ISSUER"`
	RedirectURL        string `env:"GOAP_REDIRECT_URL"`
	Host               string `env:"GOAP_HOST" envDefault:"http://localhost"`
	Port               int    `env:"GOAP_PORT" envDefault:"8080"`
	CustomTemplatePath string `env:"GOAP_CUSTOM_TEMPLATE_PATH"`
}

func Get[T any](cfg *T) error {
	return env.Parse(cfg)
}
