package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/state"
)

var version = "dev"

func main() {
	configPath := flag.String("config", "", "path to config file (optional)")
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("portwatch %s\n", version)
		os.Exit(0)
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	store, err := state.NewStore(cfg.StateFile)
	if err != nil {
		log.Fatalf("failed to initialise state store: %v", err)
	}

	notifier, err := alert.New(cfg)
	if err != nil {
		log.Fatalf("failed to initialise notifier: %v", err)
	}

	mon, err := monitor.New(cfg, store, notifier)
	if err != nil {
		log.Fatalf("failed to initialise monitor: %v", err)
	}

	log.Printf("portwatch %s starting — scanning every %s", version, cfg.ScanInterval)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := mon.Run(ctx); err != nil {
		log.Fatalf("monitor exited with error: %v", err)
	}

	log.Println("portwatch stopped")
}
