package cache

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type CachedListing struct {
	URL  string    `json:"url"`
	Date time.Time `json:"date"`
}

// LoadCache loads the cache, from the json file.
func LoadCache() (map[string]CachedListing, error) {
	scraped, err := readCache()
	if err != nil {
		return nil, fmt.Errorf("error while logging URLs: %v", err)
	}

	return scraped, nil
}

// UpdateCache writes the given map[string]CachedListing into the cache.
func UpdateCache(lastItems map[string]CachedListing) error {
	f, err := os.OpenFile("scraped.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not open scraped.json: %v", err)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)

	// Update scraped.json with new contents
	err = encoder.Encode(lastItems)
	if err != nil {
		return fmt.Errorf("Could not encode to scraped.json: %s\n", err)
	}

	return nil
}

func readCache() (map[string]CachedListing, error) {
	res := map[string]CachedListing{}

	scraped, err := os.OpenFile("scraped.json", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("could not open scraped.json: %v", err)
	}
	defer scraped.Close()

	err = json.NewDecoder(scraped).Decode(&res)
	if err != nil {
		log.Println("cannot decode cache, ignoring", err)
	}

	return res, nil
}
