package utils

import (
	"strings"
)

func lowercaseFilter(tokens []string) []string {
	for i, token := range tokens {
		tokens[i] = strings.ToLower(token)
	}
	return tokens
}

func stopwordFilter(tokens []string) []string {
	var enStopwords = map[string]struct{}{ // maybe use a stopword library, it has more words.
		"a": {}, "and": {}, "be": {}, "have": {}, "i": {},
		"in": {}, "of": {}, "that": {}, "the": {}, "to": {},
	}

	filteredTokens := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if _, ok := enStopwords[token]; !ok {
			filteredTokens = append(filteredTokens, token)
		}
	}

	return filteredTokens
}
