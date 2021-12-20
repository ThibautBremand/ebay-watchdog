package scraper

import (
	"testing"
	"time"
)

func TestParseDate(t *testing.T) {
	t.Run("US format", func(t *testing.T) {
		str := "Jun-26 06:21"
		listingURL := "https://www.ebay.com/itm/393802831789?hash=item5bb07a57ad:g:q1MAAOSwCRthvdNa"
		got, err := parseDate(str, listingURL)
		if err != nil {
			t.Errorf("error while parsing date: %v", err)
		}

		exp := time.Date(time.Now().Year(), 6, 26, 6, 21, 00, 0, time.Local)
		if exp != got {
			t.Errorf("expected %v but got %v", exp, got)
		}
	})

	t.Run("UK format", func(t *testing.T) {
		str := "26-Jun 15:39"
		listingURL := "https://www.ebay.co.uk/itm/393802831789?hash=item5bb07a57ad:g:q1MAAOSwCRthvdNa"
		got, err := parseDate(str, listingURL)
		if err != nil {
			t.Errorf("error while parsing date: %v", err)
		}

		exp := time.Date(time.Now().Year(), 6, 26, 15, 39, 00, 0, time.Local)
		if exp != got {
			t.Errorf("expected %v but got %v", exp, got)
		}
	})

	t.Run("FR format", func(t *testing.T) {
		str := "déc.-20 11:30"
		listingURL := "https://www.ebay.fr/itm/393802831789?hash=item5bb07a57ad:g:q1MAAOSwCRthvdNa"
		got, err := parseDate(str, listingURL)
		if err != nil {
			t.Errorf("error while parsing date: %v", err)
		}

		exp := time.Date(time.Now().Year(), 12, 20, 11, 30, 00, 0, time.Local)
		if exp != got {
			t.Errorf("expected %v but got %v", exp, got)
		}
	})
}

func TestParseDateByLocDomain(t *testing.T) {
	t.Run("US format", func(t *testing.T) {
		URL := "https://www.ebay.com/itm/393802831789?hash=item5bb07a57ad:g:q1MAAOSwCRthvdNa"
		date := "20 dec 2021 11:30"

		got, err := parseDateByLocDomain(date, URL)
		if err != nil {
			t.Errorf("error while parsing date: %s", err)
		}

		exp := time.Date(time.Now().Year(), 12, 20, 11, 30, 0, 0, time.Local)
		if exp != got {
			t.Errorf("expected %v but got %v", exp, got)
		}
	})

	t.Run("UK format", func(t *testing.T) {
		URL := "https://www.ebay.co.uk/itm/393802831789?hash=item5bb07a57ad:g:q1MAAOSwCRthvdNa"
		date := "20 dec 2021 11:30"

		got, err := parseDateByLocDomain(date, URL)
		if err != nil {
			t.Errorf("error while parsing date: %s", err)
		}

		exp := time.Date(time.Now().Year(), 12, 20, 11, 30, 0, 0, time.Local)
		if exp != got {
			t.Errorf("expected %v but got %v", exp, got)
		}
	})

	t.Run("FR format", func(t *testing.T) {
		URL := "https://www.ebay.fr/itm/393802831789?hash=item5bb07a57ad:g:q1MAAOSwCRthvdNa"
		date := "20 déc 2021 11:30"

		got, err := parseDateByLocDomain(date, URL)
		if err != nil {
			t.Errorf("error while parsing date: %s", err)
		}

		exp := time.Date(time.Now().Year(), 12, 20, 11, 30, 0, 0, time.Local)
		if exp != got {
			t.Errorf("expected %v but got %v", exp, got)
		}
	})
}

func TestParseLocDomain(t *testing.T) {
	URL := "https://www.ebay.fr/itm/393802831789?hash=item5bb07a57ad:g:q1MAAOSwCRthvdNa"
	got, err := parseLocDomain(URL)
	exp := "fr"

	if err != nil {
		t.Errorf("received error %s", err)
	}

	if got != exp {
		t.Errorf("expected %s but got %s", exp, got)
	}

	URL = "https://www.ebay.co.uk/itm/393802831789?hash=item5bb07a57ad:g:q1MAAOSwCRthvdNa"
	got, err = parseLocDomain(URL)
	exp = "co.uk"

	if err != nil {
		t.Errorf("received error %s", err)
	}

	if got != exp {
		t.Errorf("expected %s but got %s", exp, got)
	}
}
