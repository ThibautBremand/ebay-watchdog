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

	dayMonth := make([]string, 0)
	hours := split[1]
	if len(split) == 2 {
		dayMonth = strings.Split(split[0], "-")
		if len(dayMonth) < 2 {
			return time.Time{}, fmt.Errorf("error while parsing date %s", str)
		}
	} else if len(split) == 3 {
		// When the day and month are not separated by a '-' but by a space.
		dayMonth = append(dayMonth, split[0])
		dayMonth = append(dayMonth, split[1])
		hours = split[2]
	} else {
		return time.Time{}, fmt.Errorf("error: date %s is in unknown format and cannot be parsed", str)
	}

	day := dayMonth[0]
	month := dayMonth[1]
	if len(day) > len(month) {
		temp := day
		day = month
		month = temp
	}

	if len(month) > 3 {
		month = firstN(month, 3)
	}

	if len(day) > 2 {
		day = firstN(day, 2)
	}

	year, _, _ := time.Now().Date()
	fullDate := fmt.Sprintf("%s %s %d %s", day, month, year, hours)

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
	case "com", "ca", "com.au", "com.sg", "com.my", "ph":
		locale = monday.LocaleEnUS
	case "co.uk", "ie":
		locale = monday.LocaleEnGB
	case "fr":
		locale = monday.LocaleFrFR
	case "de", "ch", "at":
		locale = monday.LocaleDeDE
	case "es":
		locale = monday.LocaleEsES
	case "it":
		locale = monday.LocaleItIT
	case "nl":
		locale = monday.LocaleNlNL
	case "pl":
		locale = monday.LocalePlPL
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

// firstN returns the first n characters of a string, and it correctly counts the unicode characters as 1.
func firstN(s string, n int) string {
	i := 0
	for j := range s {
		if i == n {
			return s[:j]
		}
		i++
	}
	return s
}
