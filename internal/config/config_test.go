package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("HOST", "")
	t.Setenv("PORT", "")
	t.Setenv("ENV", "")
	t.Setenv("STORAGE_TYPE", "")

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
}

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("HOST", " 0.0.0.0 ")
	t.Setenv("PORT", " 9090 ")
	t.Setenv("ENV", " dev ")
	t.Setenv("STORAGE_TYPE", " memory ")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Host != "0.0.0.0" {
		t.Fatalf("Host = %q, want %q", cfg.Host, "0.0.0.0")
	}
	if cfg.Port != "9090" {
		t.Fatalf("Port = %q, want %q", cfg.Port, "9090")
	}
	if cfg.Env != "dev" {
		t.Fatalf("Env = %q, want %q", cfg.Env, "dev")
	}
	if cfg.AddrString() != "0.0.0.0:9090" {
		t.Fatalf("AddrString() = %q, want %q", cfg.AddrString(), "0.0.0.0:9090")
	}
	if !cfg.IsMemoryStorage() {
		t.Fatalf("IsMemoryStorage() = false, want true")
	}
}
