package scraper

import (
	"fmt"
	"github.com/goodsign/monday"
	"strings"
	"time"
)

// parseDate returns a time.Time from the given date as string. It handles multiple language formats determined from
// the given listing URL.
// For example, '.com' will be handled as US english, '.co.uk' as UK english, '.fr' as french.
//
// US example: Jun-26 06:21
// UK example: 26-Jun 15:39
// FR example: déc.-20 11:30
func parseDate(str string, URL string) (time.Time, error) {
	split := strings.Split(str, " ")
	if len(split) < 2 {
		return time.Time{}, fmt.Errorf("error while parsing date %s", str)
	}

	first := strings.Split(split[0], "-")
	if len(first) < 2 {
		return time.Time{}, fmt.Errorf("error while parsing date %s", str)
	}

	day := first[0]
	month := first[1]
	if len(day) > len(month) {
		temp := day
		day = month
		month = temp
	}

	if len(month) > 3 {
		month = month[:4]
	}

	year, _, _ := time.Now().Date()
	fullDate := fmt.Sprintf("%s %s %d %s", day, month, year, split[1])

	t, err := parseDateByLocDomain(fullDate, URL)
	if err != nil {
		return time.Time{}, fmt.Errorf("error while parsing date %s: %v", str, err)
	}

	return t, nil
}

// parseDateByLocDomain returns the time.Time from the given date as string and the given listing URL.
// The listing URL is used in order to detect and corresponding location (US, UK, FR, etc) and parse the date
// accordingly.
func parseDateByLocDomain(date string, URL string) (time.Time, error) {
	locDomain, err := parseLocDomain(URL)
	if err != nil {
		return time.Time{}, fmt.Errorf("could not parse location domain from URL %s: %s", URL, err)
	}

	var locale monday.Locale
	switch locDomain {
	case "com":
		locale = monday.LocaleEnUS
	case "co.uk":
		locale = monday.LocaleEnGB
	case "fr":
		locale = monday.LocaleFrFR
	default:
		return time.Time{}, fmt.Errorf("unhandled loc domain: %s", locDomain)
	}

	loc, err := time.LoadLocation("Local")
	if err != nil {
		return time.Time{}, fmt.Errorf("could not load location: %s", err)
	}

	t, err := monday.ParseInLocation("2 Jan 2006 15:04", date, loc, locale)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
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
