package utils

import (
	"strings"
	"unicode"
)

func tokenize(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsNumber(r) && !unicode.IsLetter(r)
	})
}

func analyze(text string) []string {
	tokens := tokenize(text)
	tokens = lowercaseFilter(tokens)
	// add an extra step to detect/identify language which can be used in stopwords filter and stemming, there is some effort in this area
	tokens = stopwordFilter(tokens)
	tokens = stemmingTransform(tokens)
	return tokens
}
