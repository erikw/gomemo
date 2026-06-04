package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Host string
	Port string
}

func Load() (Config, error) {
	host := strings.TrimSpace(os.Getenv("HOST"))
	if host == "" {
		host = "127.0.0.1"
	}

	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8080"
	}

	return Config{
		Port: port,
		Host: host,
	}, nil
}

func (cfg Config) String() string {
	type ConfigStringer Config
	return fmt.Sprintf("Config: %+v", ConfigStringer(cfg))
}

func (cfg Config) AddrString() string {
	return fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
}
