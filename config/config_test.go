package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	// Unset env vars to test defaults
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("SERVER_MODE")

	cfg := Load()

	if cfg.Database.Host != "localhost" {
		t.Errorf("expected DB_HOST default 'localhost', got '%s'", cfg.Database.Host)
	}
	if cfg.Database.Port != "5432" {
		t.Errorf("expected DB_PORT default '5432', got '%s'", cfg.Database.Port)
	}
	if cfg.Database.User != "postgres" {
		t.Errorf("expected DB_USER default 'postgres', got '%s'", cfg.Database.User)
	}
	if cfg.Database.Password != "password" {
		t.Errorf("expected DB_PASSWORD default 'password', got '%s'", cfg.Database.Password)
	}
	if cfg.Database.Name != "budget_tracker" {
		t.Errorf("expected DB_NAME default 'budget_tracker', got '%s'", cfg.Database.Name)
	}
	if cfg.JWT.Secret != "default-secret-change-in-production" {
		t.Errorf("expected JWT_SECRET default, got '%s'", cfg.JWT.Secret)
	}
	if cfg.Server.Port != "8080" {
		t.Errorf("expected SERVER_PORT default '8080', got '%s'", cfg.Server.Port)
	}
	if cfg.Server.Mode != "debug" {
		t.Errorf("expected SERVER_MODE default 'debug', got '%s'", cfg.Server.Mode)
	}
}

func TestLoad_EnvOverride(t *testing.T) {
	os.Setenv("DB_HOST", "remotehost")
	os.Setenv("DB_PORT", "6543")
	os.Setenv("DB_USER", "myuser")
	os.Setenv("DB_PASSWORD", "mypassword")
	os.Setenv("DB_NAME", "mydb")
	os.Setenv("JWT_SECRET", "supersecret")
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("SERVER_MODE", "release")

	cfg := Load()

	if cfg.Database.Host != "remotehost" {
		t.Errorf("expected DB_HOST 'remotehost', got '%s'", cfg.Database.Host)
	}
	if cfg.Database.Port != "6543" {
		t.Errorf("expected DB_PORT '6543', got '%s'", cfg.Database.Port)
	}
	if cfg.Database.User != "myuser" {
		t.Errorf("expected DB_USER 'myuser', got '%s'", cfg.Database.User)
	}
	if cfg.Database.Password != "mypassword" {
		t.Errorf("expected DB_PASSWORD 'mypassword', got '%s'", cfg.Database.Password)
	}
	if cfg.Database.Name != "mydb" {
		t.Errorf("expected DB_NAME 'mydb', got '%s'", cfg.Database.Name)
	}
	if cfg.JWT.Secret != "supersecret" {
		t.Errorf("expected JWT_SECRET 'supersecret', got '%s'", cfg.JWT.Secret)
	}
	if cfg.Server.Port != "9090" {
		t.Errorf("expected SERVER_PORT '9090', got '%s'", cfg.Server.Port)
	}
	if cfg.Server.Mode != "release" {
		t.Errorf("expected SERVER_MODE 'release', got '%s'", cfg.Server.Mode)
	}
}

func TestGetEnv(t *testing.T) {
	os.Setenv("TEST_KEY", "value")
	if v := getEnv("TEST_KEY", "default"); v != "value" {
		t.Errorf("expected 'value', got '%s'", v)
	}
	os.Unsetenv("TEST_KEY")
	if v := getEnv("TEST_KEY", "default"); v != "default" {
		t.Errorf("expected 'default', got '%s'", v)
	}
}

func TestGetEnvAsInt(t *testing.T) {
	os.Setenv("TEST_INT", "123")
	if v := getEnvAsInt("TEST_INT", 42); v != 123 {
		t.Errorf("expected 123, got %d", v)
	}
	os.Setenv("TEST_INT", "notanint")
	if v := getEnvAsInt("TEST_INT", 42); v != 42 {
		t.Errorf("expected fallback 42, got %d", v)
	}
	os.Unsetenv("TEST_INT")
	if v := getEnvAsInt("TEST_INT", 42); v != 42 {
		t.Errorf("expected fallback 42, got %d", v)
	}
}
