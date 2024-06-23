package utils

import (
	"github.com/kljensen/snowball"
)

func stemmingTransform(tokens []string) []string {
	for i, token := range tokens {
		stemmedToken, err := snowball.Stem(token, "english", false)
		if err != nil {
			tokens[i] = stemmedToken
		}
	}
	return tokens
}
