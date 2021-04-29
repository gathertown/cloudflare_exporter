package main

import (
	"os"

	config "github.com/gathertown/cloudflare_exporter/internal/config"
	log "github.com/gathertown/cloudflare_exporter/pkg/log"
)

// run labels nodes if label is missing
func main() {
	cfg := config.FromEnv()
	logger := log.New(os.Stdout, cfg.Env)
	logger.Info("Launching cloudformation_exporter", "zoneTag", cfg.ZoneTag)
}
