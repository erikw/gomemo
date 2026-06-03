package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port string
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return Config{
		Port: port,
	}
}

func (cfg Config) String() string {
	type ConfigStringer Config
	return fmt.Sprintf("Config: %+v", ConfigStringer(cfg))
}

func (cfg Config) AddrString() string {
	return fmt.Sprintf(":%s", cfg.Port)
}
