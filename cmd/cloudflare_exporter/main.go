package main

import (
	"fmt"
	"os"

	config "github.com/gathertown/cloudflare_exporter/internal/config"
	r "github.com/gathertown/cloudflare_exporter/pkg/cloudflare/requests"
	log "github.com/gathertown/cloudflare_exporter/pkg/log"
)

// run labels nodes if label is missing
func main() {
	cfg := config.FromEnv()
	logger := log.New(os.Stdout, cfg.Env)
	logger.Info("Launching cloudformation_exporter", "zoneTag", cfg.ZoneTag)
	a, err := r.Requests(cfg.ZoneTag, cfg.Token)
	if err != nil {
		panic(err)
	}
	fmt.Println(a)
}
