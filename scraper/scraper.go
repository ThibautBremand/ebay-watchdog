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

	log.Println("Got new listings!", len(listings))
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

	for _, searchItem := range searchItems {
		searchUrl := searchItem.URL

		log.Println("Searching with", searchUrl)
		doc, err := web.Get(searchUrl)
		if err != nil {
			log.Println("could not make request to SearchItem page", err)
			continue
		}

		isFirst := true

		doc.Find("div.s-item__info").EachWithBreak(func(i int, sel *goquery.Selection) bool {
			listing, b := parseItem(sel, scraped, searchUrl)
			if listing != nil {
				pulledListings = append(pulledListings, *listing)
				if isFirst {
					lastItems[searchUrl] = *listing
					isFirst = false
				}
			}

			return b
		})
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

	url, exists := itemSel.Attr("href")
	if !exists {
		return nil, true
	}

	if isKnownURL && scraped[searchUrl].URL == url {
		log.Println("Found nothing new")
		return nil, false
	}

	title := sel.Find(".s-item__title").Text()
	subtitle := sel.Find(".s-item__subtitle").Text()

	detailsSel := sel.Find(".s-item__details").Children()

	price := detailsSel.Find(".s-item__price").Text()
	date := detailsSel.Find(".s-item__listingDate").Text()

	t, err := parseDate(date)
	if err != nil {
		log.Println("error while parsing date", date, err)
		return nil, false
	}

	lastScrapedProductDate := scraped[searchUrl].Date.Add(time.Hour * time.Duration(-1))
	// In case the last scraped product has been deleted, we can still compare the dates
	if isKnownURL && t.Before(lastScrapedProductDate) {
		log.Println("Found nothing new")
		return nil, false
	}

	listing := Listing{
		URL:      url,
		Title:    title,
		Subtitle: subtitle,
		Price:    price,
		Date:     t,
	}

	log.Println("Got listing details")

	if !isKnownURL {
		// This was the first time scraping this searchUrl. As we only want to check for new listings,
		// we won't scrape all the next listings and we will just wait for new ones. This is why we
		// will break out of the loop.
		return &listing, false
	}

	return &listing, true
}

// parseDate returns a time.Time from the given date as string. It handles dates from eBay listings in US english and
// UK english format.
//
// US example: Jun-26 06:21
// UK example: 26-Jun 15:39
func parseDate(str string) (time.Time, error) {
	split := strings.Split(str, " ")
	if len(split) < 2 {
		return time.Time{}, fmt.Errorf("error while parsing date %s", str)
	}

	first := strings.Split(split[0], "-")
	if len(first) < 2 {
		return time.Time{}, fmt.Errorf("error while parsing date %s", str)
	}

	reordered := fmt.Sprintf("%s %s", first[0], first[1])
	if len(first[0]) == 3 {
		reordered = fmt.Sprintf("%s %s", first[1], first[0])
	}

	year, _, _ := time.Now().Date()
	fullDate := fmt.Sprintf("%s %d %s", reordered, year, split[1])

	t, err := time.Parse("2 Jan 2006 15:04", fullDate)
	if err != nil {
		return time.Time{}, fmt.Errorf("error while parsing date %s: %v", str, err)
	}

	return t, nil
}
