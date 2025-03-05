package utils

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func NormalizeCityName(name string) string {
	caser := cases.Title(language.Russian)
	return caser.String(strings.ToLower(name))
}
