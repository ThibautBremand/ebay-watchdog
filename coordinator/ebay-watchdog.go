package coordinator

import (
	"bytes"
	"ebay-watchdog/cache"
	"ebay-watchdog/config"
	"ebay-watchdog/scraper"
	"ebay-watchdog/web"
	"github.com/spf13/viper"
	"log"
	"strings"
	"text/template"
	"time"
)

func Start(
	searchItems []config.SearchItem,
	tpl *template.Template,
	scrapedURLs map[string]cache.CachedListing,
	sleepPeriod time.Duration,
) {
	for {
		listings, lastItems, err := scraper.Scrape(searchItems, scrapedURLs)
		if err != nil {
			log.Println("error while scraping new listings, skipping", err)
			time.Sleep(sleepPeriod)
			continue
		}

		scrapedURLs = buildCache(lastItems, scrapedURLs)

		err = cache.UpdateCache(scrapedURLs)
		if err != nil {
			log.Println("error while updating scraped URLs, skipping", err)
		}

		sendToTelegram(listings, tpl)

		time.Sleep(sleepPeriod)
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
			viper.GetString("TELEGRAM_TOKEN"),
			viper.GetString("TELEGRAM_CHAT_ID"),
			msg,
		)
		if err != nil {
			log.Println("could not send Telegram message", err)
		}
	}
}
