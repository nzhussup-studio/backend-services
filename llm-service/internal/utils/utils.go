package utils

import "time"

func IsValidLanguage(lang string) bool {
	validLanguages := map[string]bool{
		"en": true,
		"kk": true,
		"de": true,
	}

	return validLanguages[lang]
}

func ToPtrTime(t time.Time) *time.Time {
	return &t
}
