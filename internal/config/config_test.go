package config

import (
	"os"
	"testing"
)

func setenv(t *testing.T, key, value string) {
	t.Helper()
	t.Logf("Setting env %q=%q", key, value)
	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("Failed setting env %q as %q: %v", key, value, err)
	}
}

func unsetenv(t *testing.T, key string) {
	t.Helper()
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("Failed unsetting env %q: %v", key, err)
	}
}

func TestFromEnv(t *testing.T) {
	setenv(t, "ENV", "development")
	setenv(t, "TOKEN", "a1234")
	setenv(t, "SUBSYSTEM", "my_org")
	setenv(t, "ZONETAG", "abcde")
	setenv(t, "ORIGIN", "origin")

	cfg := FromEnv()

	if got, want := cfg.Env, "development"; got != want {
		t.Errorf("FromEnv() 'ENV' = %q; want %q", got, want)
	}

	if got, want := cfg.Token, "a1234"; got != want {
		t.Errorf("FromEnv() 'TOKEN' = %q; want %q", got, want)
	}

	if got, want := cfg.ZoneTag, "abcde"; got != want {
		t.Errorf("FromEnv() 'ZONETAG' = %q; want %q", got, want)
	}

	if got, want := cfg.Sub, "my_org"; got != want {
		t.Errorf("FromEnv() 'SUBSYSTEM' = %q; want %q", got, want)
	}

	if got, want := cfg.Origin, "origin"; got != want {
		t.Errorf("FromEnv() 'ORIGIN' = %q; want %q", got, want)
	}

	unsetenv(t, "ENV")
	unsetenv(t, "TOKEN")
	unsetenv(t, "ZONE")
	unsetenv(t, "SUBSYSTEM")
	unsetenv(t, "ORIGIN")
}
