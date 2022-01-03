package scraper

import (
	"ebay-watchdog/cache"
	"ebay-watchdog/config"
	"ebay-watchdog/web"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strings"
	"time"
)

type Listing struct {
	URL      string    `json:"url"`
	Title    string    `json:"title"`
	Subtitle string    `json:"subtitle"`
	Price    string    `json:"price"`
	Date     time.Time `json:"date"`
	ID       string    `json:"id"`
}

// Scrape starts the scraping for the given []config.SearchItem.
// It returns a list of Listing, to be sent to Telegram. It also returns a map[string]Listing which will be used when
// updating the cache.
func Scrape(
	searchItems []config.SearchItem,
	cache map[string]cache.CachedListing,
) (
	[]Listing,
	map[string]Listing,
	error,
) {
	log.Println("Scraping new listings")
	listings, lastItems, err := scrapeListings(searchItems, cache)
	if err != nil {
		return nil, nil, fmt.Errorf("could not start scraping: %v", err)
	}

	log.Printf("Got %d new listings!\n", len(listings))
	return listings, lastItems, nil
}

func scrapeListings(
	searchItems []config.SearchItem,
	scraped map[string]cache.CachedListing,
) (
	[]Listing,
	map[string]Listing,
	error,
) {
	var pulledListings []Listing
	lastItems := make(map[string]Listing)

	// Keep in memory the id of the parsed listings, so we do not send the same listing twice when checking
	// multiple domains.
	currentSearchItems := make(map[string]int)

	for _, searchItem := range searchItems {
		if searchItem.Domains == nil || len(searchItem.Domains) == 0 {
			domain, err := parseLocDomain(searchItem.URL)
			if err != nil {
				return nil, nil, fmt.Errorf("could not get domain from url %s: %s", searchItem.URL, err)
			}

			searchItem.Domains = []string{domain}
		}

		for _, domain := range searchItem.Domains {
			searchUrl, err := setDomain(searchItem.URL, domain)
			if err != nil {
				return nil, nil, fmt.Errorf("could not set domain %s for url %s: %s", domain, searchItem.URL, err)
			}

			log.Printf("Searching with url %s (domain %s)\n", searchUrl, domain)
			doc, err := web.Get(searchUrl)
			if err != nil {
				log.Println("could not make request to SearchItem page", err)
				continue
			}

			isFirst := true

			doc.Find("div.s-item__info").EachWithBreak(func(i int, sel *goquery.Selection) bool {
				listing, b := parseItem(sel, scraped, searchUrl)
				if listing != nil {
					_, isKnownID := currentSearchItems[listing.ID]
					if !isKnownID {
						currentSearchItems[listing.ID] = 1
						pulledListings = append(pulledListings, *listing)
					}

					if isFirst {
						lastItems[searchUrl] = *listing

						isFirst = false
					}
				}

				return b
			})

			// We space each queries just in case, to prevent getting throttled
			time.Sleep(2 * time.Second)
		}
	}

	return pulledListings, lastItems, nil
}

func parseItem(
	sel *goquery.Selection,
	scraped map[string]cache.CachedListing,
	searchUrl string,
) (*Listing, bool) {
	_, isKnownURL := scraped[searchUrl]

	itemSel := sel.Children()
	if len(itemSel.Nodes) < 3 {
		return nil, true
	}

	rawURL, exists := itemSel.Attr("href")
	if !exists {
		return nil, true
	}

	// Listing URLs with amdata generate different URLs for the same listings
	// Removing the amdata allows us to determine if a listing has already been scraped or not.
	URL := strings.Split(rawURL, "&amdata")[0]

	if isKnownURL && scraped[searchUrl].URL == URL {
		log.Println("Stop - Reached a listing that has already been scraped!")
		return nil, false
	}

	title := sel.Find(".s-item__title").Text()
	subtitle := sel.Find(".s-item__subtitle").Text()
	detailsSel := sel.Find(".s-item__details").Children()
	price := detailsSel.Find(".s-item__price").Text()
	date := detailsSel.Find(".s-item__listingDate").Text()

	t, err := parseDate(date, URL)
	if err != nil {
		log.Println("error while parsing date", date, err)
		return nil, false
	}

	lastScrapedProductDate := scraped[searchUrl].Date.Add(time.Hour * time.Duration(-1))
	// In case the last scraped product has been deleted, we can still compare the dates
	if isKnownURL && t.Before(lastScrapedProductDate) {
		log.Println("Stop - Reached a listing that has an older publication date than the last scraped listing!")
		return nil, false
	}

	split := strings.Split(URL, "/")

	listing := Listing{
		URL:      URL,
		Title:    title,
		Subtitle: subtitle,
		Price:    price,
		Date:     t,
		ID:       split[len(split)-1],
	}

	log.Printf("Successfully scraped 1 listing details (ID: %s)\n", listing.ID)

	if !isKnownURL {
		// This was the first time scraping this searchUrl. As we only want to check for new listings,
		// we won't scrape all the next listings and we will just wait for new ones. This is why we
		// will break out of the loop.
		return &listing, false
	}

	return &listing, true
}

// setDomain replaces the top level domain of the given URL by the given domain, and returns the new URL.
func setDomain(URL string, domain string) (string, error) {
	split := strings.Split(URL, "/")

	if len(split) < 3 {
		return "", fmt.Errorf("could not extract domain from URL %s", URL)
	}

	split[2] = fmt.Sprintf("%s%s", "www.ebay.", domain)

	return strings.Join(split, "/"), nil
}

// parseLocDomain returns the location domain from the given listing URL.
// e.g. com, co.uk, fr, etc.
func parseLocDomain(URL string) (string, error) {
	split := strings.Split(URL, "/")

	if len(split) < 3 {
		return "", fmt.Errorf("could not extract location domain from URL %s", URL)
	}

	return strings.Replace(split[2], "www.ebay.", "", 1), nil
}
