package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type ZohoConfig struct {
	ClientID     string `toml:"client_id"`
	ClientSecret string `toml:"client_secret"`
	RefreshToken string `toml:"refresh_token"`
}

type Config struct {
	ExpectedToken string     `toml:"expected_token"`
	Zoho          ZohoConfig `toml:"zoho"`
}

var Cfg Config

func LoadConfig(path string) error {
	if _, err := toml.DecodeFile(path, &Cfg); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	return nil
}
