package index

import (
	"strings"
	"unicode/utf8"
)

const minWordLength = 2

var InvertedIndexMap = make(map[string][]string)

func InvertIndex(inputWords []string, docId string) {

	inputWords = cleanText(inputWords)
	for _, word := range inputWords {
		if !isStringInSlice(docId, InvertedIndexMap[word]) {
			InvertedIndexMap[word] = append(InvertedIndexMap[word], docId)
		}
	}
}

func cleanText(inputWords []string) []string {

	cleanWords := make([]string, 0)
	for _, word := range inputWords {

		if stopWORDS[word] || utf8.RuneCountInString(word) < minWordLength {
			continue
		}
		word = strings.ToLower(word)
		cleanWords = append(cleanWords, word)
	}
	return cleanWords
}

func isStringInSlice(a string, list []string) bool {

	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
