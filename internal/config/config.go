package config

import (
	"fmt"
	"os"
	"strings"
)

const defaultHost = "127.0.0.1"
const defaultPort = "8080"
const defaultEnv = "prod"
const defaultStorageType = "memory"

type Config struct {
	Host        string
	Port        string
	Env         string
	StorageType string
}

func Load() (Config, error) {
	host := strings.TrimSpace(os.Getenv("HOST"))
	if host == "" {
		host = defaultHost
	}

	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = defaultPort
	}

	env := strings.TrimSpace(os.Getenv("ENV"))
	if env == "" {
		env = defaultEnv
	}

	storageType := strings.TrimSpace(os.Getenv("STORAGE_TYPE"))
	if storageType == "" {
		storageType = defaultStorageType
	}

	return Config{
		Host:        host,
		Port:        port,
		Env:         env,
		StorageType: storageType,
	}, nil
}

func (cfg Config) String() string {
	type ConfigStringer Config
	return fmt.Sprintf("Config: %+v", ConfigStringer(cfg))
}

func (cfg Config) AddrString() string {
	return fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
}

func (cfg Config) IsMemoryStorage() bool {
	return cfg.StorageType == "memory"
}
