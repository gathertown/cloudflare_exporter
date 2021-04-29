package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gathertown/cloudflare_exporter/internal/config"
)

// run labels nodes if label is missing
func main() {
	cfg := config.FromEnv()
	logger := log.New(os.Stdout, cfg.Env)
	logger.Info("Launching casper-3", "labelKey", cfg.LabelKey, "labelValue", cfg.LabelValue, "interval", cfg.ScanIntervalSeconds, "environment", cfg.Env, "TXT identifier", fmt.Sprintf("heritage=casper-3,environment=%s", cfg.Env))
}
