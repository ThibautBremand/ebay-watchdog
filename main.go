package main

import (
	"ebay-watchdog/cache"
	"ebay-watchdog/config"
	"ebay-watchdog/coordinator"
	"log"
	"time"
)

func main() {
	log.Println("Starting ebay-watchdog")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("error while loading config: %v", err)
	}

	tpl, err := cfg.LoadTemplate()
	if err != nil {
		log.Fatalf("Could not parse message template: %v\n", err)
	}

	scrapedURLs, err := cache.LoadCache()
	if err != nil {
		log.Fatalf("Could not load scraper urls: %v", err)
	}

	sleepPeriod := time.Duration(cfg.Delay) * time.Second

	c := coordinator.NewCoordinator(cfg.Searches, sleepPeriod, tpl)
	c.Start(scrapedURLs)
}
