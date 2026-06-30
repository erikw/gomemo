package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("HOST", "")
	t.Setenv("PORT", "")
	t.Setenv("ENV", "")
	t.Setenv("STORAGE_TYPE", "")
	t.Setenv("DATABASE_URL", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Host != defaultHost {
		t.Fatalf("Host = %q, want %q", cfg.Host, defaultHost)
	}
	if cfg.Port != defaultPort {
		t.Fatalf("Port = %q, want %q", cfg.Port, defaultPort)
	}
	if cfg.Env != defaultEnv {
		t.Fatalf("Env = %q, want %q", cfg.Env, defaultEnv)
	}
	if cfg.StorageType != defaultStorageType {
		t.Fatalf("StorageType = %q, want %q", cfg.StorageType, defaultStorageType)
	}
	if cfg.DatabaseURL != "" {
		t.Fatalf("DatabaseURL = %q, want empty", cfg.DatabaseURL)
	}
}

func TestLoadRequiresDatabaseURLForPostgres(t *testing.T) {
	t.Setenv("STORAGE_TYPE", "postgres")
	t.Setenv("DATABASE_URL", "")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() expected error for missing DATABASE_URL")
	}
}

func TestLoadRejectsUnsupportedStorageType(t *testing.T) {
	t.Setenv("STORAGE_TYPE", "sqlite")
	t.Setenv("DATABASE_URL", "")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() expected error for unsupported STORAGE_TYPE")
	}
}
