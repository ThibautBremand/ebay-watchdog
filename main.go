package main

import (
	"ebay-watchdog/cache"
	"ebay-watchdog/config"
	"ebay-watchdog/coordinator"
	"github.com/spf13/viper"
	"log"
	"time"
)

func main() {
	log.Println("Starting ebay-watchdog")

	err := config.Load()
	if err != nil {
		log.Fatalf("error while loading config: %v", err)
	}

	searchItems, err := config.LoadSearchItems()
	if err != nil {
		log.Fatalf("Could not get SearchItem URLs: %v\n", err)
	}

	tpl, err := config.LoadTemplate()
	if err != nil {
		log.Fatalf("Could not parse message template: %v\n", err)
	}

	scrapedURLs, err := cache.LoadCache(viper.GetBool("track-scraped-urls"))
	if err != nil {
		log.Fatalf("Could not load scraper urls: %v", err)
	}

	sleepPeriod := time.Duration(viper.GetInt("delay")) * time.Second

	coordinator.Start(searchItems, tpl, scrapedURLs, sleepPeriod)
}
