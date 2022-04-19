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

		str = "avr.-17 21:08"
		listingURL = "https://www.ebay.fr/itm/393802831789?hash=item5bb07a57ad:g:q1MAAOSwCRthvdNa"
		got, err = parseDate(str, listingURL)
		if err != nil {
			t.Errorf("error while parsing date: %v", err)
		}

		exp = time.Date(time.Now().Year(), 4, 17, 21, 8, 00, 0, time.Local)
		if exp != got {
			t.Errorf("expected %v but got %v", exp, got)
		}
	})

	t.Run("DE format", func(t *testing.T) {
		str := "17. Nov. 04:09"
		listingURL := "https://www.ebay.de/itm/393802831789?hash=item5bb07a57ad:g:q1MAAOSwCRthvdNa"
		got, err := parseDate(str, listingURL)
		if err != nil {
			t.Errorf("error while parsing date: %v", err)
		}

		exp := time.Date(time.Now().Year(), 11, 17, 4, 9, 00, 0, time.Local)
		if exp != got {
			t.Errorf("expected %v but got %v", exp, got)
		}
	})
}

func TestFirstN(t *testing.T) {
	t.Run("Unicode", func(t *testing.T) {
		got := firstN("世界 Hello", 1)
		exp := "世"

		if exp != got {
			t.Errorf("expected %s but got %s", exp, got)
		}
	})

	t.Run("Non unicode", func(t *testing.T) {
		got := firstN("Hello World", 1)
		exp := "H"

		if exp != got {
			t.Errorf("expected %s but got %s", exp, got)
		}
	})
}
