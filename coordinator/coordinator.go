package coordinator

import (
	"bytes"
	"ebay-watchdog/cache"
	"ebay-watchdog/config"
	"ebay-watchdog/scraper"
	"ebay-watchdog/web"
	"log"
	"os"
	"strings"
	"text/template"
	"time"
)

type Coordinator struct {
	Scraper     *scraper.Scraper
	SleepPeriod time.Duration
	Tpl         *template.Template
}

func NewCoordinator(
	searchItems []config.SearchItem,
	sleepPeriod time.Duration,
	tpl *template.Template,
) *Coordinator {
	searchURLs := buildSearchURLs(searchItems)
	s := scraper.NewScraper(searchURLs)

	return &Coordinator{
		Scraper:     s,
		SleepPeriod: sleepPeriod,
		Tpl:         tpl,
	}
}

func (c *Coordinator) Start(
	scrapedURLs map[string]cache.CachedListing,
) {
	for {
		listings, lastItems, err := c.Scraper.Scrape(scrapedURLs)
		if err != nil {
			log.Println("error while scraping new listings, skipping", err)
			time.Sleep(c.SleepPeriod)
			continue
		}

		scrapedURLs = buildCache(lastItems, scrapedURLs)

		err = cache.UpdateCache(scrapedURLs)
		if err != nil {
			log.Println("error while updating scraped URLs, skipping", err)
		}

		sendToTelegram(listings, c.Tpl)

		time.Sleep(c.SleepPeriod)
	}

}

// buildCache returns a map[string]cache.CachedListing, ready to be persisted into the cache, from the given
// map[string]scraper.Listing which comes from the last scraping, and the map[string]cache.CachedListing which is the
// previous cache.
// It uses data from both maps to build the new cache.
func buildCache(lastItems map[string]scraper.Listing, scrapedURLs map[string]cache.CachedListing) map[string]cache.CachedListing {
	lastScrapedURLs := make(map[string]cache.CachedListing)
	for key, listing := range lastItems {
		toPersist := cache.CachedListing{
			URL:  listing.URL,
			Date: listing.Date,
		}

		lastScrapedURLs[key] = toPersist
	}

	for k, v := range scrapedURLs {
		if _, ok := lastScrapedURLs[k]; ok {
			continue
		}
		lastScrapedURLs[k] = v
	}
	return lastScrapedURLs
}

func sendToTelegram(listings []scraper.Listing, tpl *template.Template) {
	for _, listing := range listings {
		buf := &bytes.Buffer{}
		err := tpl.Execute(buf, listing)
		var msg string
		if err != nil {
			log.Println("could not execute template", err)
			msg = listing.URL
		} else {
			msg = buf.String()
		}

		// Double quotes are not correctly parsed by Telegram
		msg = strings.ReplaceAll(msg, `"`, "")

		err = web.SendTelegramMessage(
			os.Getenv("TELEGRAM_TOKEN"),
			os.Getenv("TELEGRAM_CHAT_ID"),
			msg,
		)
		if err != nil {
			log.Println("could not send Telegram message", err)
		}
	}
}

// buildSearchURLs takes a list []config.SearchItem from the config, and returns a list []scraper.SearchURL directly
// usable by the scraper.
func buildSearchURLs(searchItems []config.SearchItem) []scraper.SearchURL {
	searchURLs := make([]scraper.SearchURL, len(searchItems))
	for i, s := range searchItems {
		searchURLs[i] = scraper.SearchURL{
			URL:     s.URL,
			Domains: s.Domains,
		}
	}

	return searchURLs
}
