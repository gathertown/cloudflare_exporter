// Package config provides functions that allow to construct the service
// configuration from the environment.
package config

import (
	"os"
)

const (
	defaultEnv     = "development"
	defaultToken   = "a1234"
	defaultZoneTag = "abcde"
)

// Config contains service information that can be changed from the
// environment.
type Config struct {
	Env     string
	Token   string
	ZoneTag string
}

// FromEnv returns the service configuration from the environment variables.
// If an environment variable is not found, then a default value is provided.
func FromEnv() *Config {
	var (
		env     = getenv("ENV", defaultEnv)
		token   = getenv("TOKEN", defaultToken)
		zoneTag = getenv("ZONETAG", defaultZoneTag)
	)

	c := &Config{
		Env:     env,
		Token:   token,
		ZoneTag: zoneTag,
	}
	return c
}

func getenv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
