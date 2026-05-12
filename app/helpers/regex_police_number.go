package helpers

import (
	"regexp"
	"strings"
)

var PoliceNumberRegex = regexp.MustCompile(
	`^(?:[A-Z]{1,2})\s?[0-9]{1,4}\s?(?:[A-Z]{1,3})?$`,
)

func NormalizePoliceNumber(police string) string {
	police = strings.TrimSpace(police)
	police = strings.ToUpper(police)

	police = regexp.MustCompile(`\s+`).ReplaceAllString(police, " ")

	return police
}
