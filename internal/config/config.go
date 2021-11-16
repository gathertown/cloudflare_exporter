// Package config provides functions that allow to construct the service
// configuration from the environment.
package config

import (
	"os"
)

const (
	defaultEnv     = "development"
	defaultSub     = "gather_town"
	defaultToken   = "a1234"
	defaultZoneTag = "abcd123"
	defaultOrigin  = ""
)

// Config contains service information that can be changed from the
// environment.
type Config struct {
	Env     string
	Sub     string
	Token   string
	ZoneTag string
	Origin  string
}

// FromEnv returns the service configuration from the environment variables.
// If an environment variable is not found, then a default value is provided.
func FromEnv() *Config {
	var (
		env     = getenv("ENV", defaultEnv)
		sub     = getenv("SUBSYSTEM", defaultSub)
		token   = getenv("TOKEN", defaultToken)
		zoneTag = getenv("ZONETAG", defaultZoneTag)
		origin  = getenv("ORIGIN", defaultOrigin)
	)

	c := &Config{
		Env:     env,
		Sub:     sub,
		Token:   token,
		ZoneTag: zoneTag,
		Origin:  origin,
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
