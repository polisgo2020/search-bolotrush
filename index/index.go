package index

import (
	"strings"
	"unicode/utf8"
)

const minWordLength = 2

type wordStruct struct {
	Doc      string
	Position []int
}
type InvMap map[string][]wordStruct

var InvertedIndexMap = make(InvMap)

func (p InvMap) isWordInList(word string, docId string) (int, bool) {
	for i, ind := range p[word] {
		if ind.Doc == docId {
			return i, true
		}
	}
	return -1, false
}

func InvertIndex(inputWords []string, docId string) {

	inputWords = cleanText(inputWords)
	for i, word := range inputWords {
		if index, ok := InvertedIndexMap.isWordInList(word, docId); !ok {

			structure := wordStruct{
				Doc:      docId,
				Position: make([]int, 0),
			}

			structure.Position = append(structure.Position, i)
			InvertedIndexMap[word] = append(InvertedIndexMap[word], structure)
		} else if index != -1 {
			InvertedIndexMap[word][index].Position = append(InvertedIndexMap[word][index].Position, i)
		}
	}
}

func GetDocStrSlice(slice []wordStruct) []string {
	outSlice := make([]string, 0)
	for _, doc := range slice {
		outSlice = append(outSlice, doc.Doc)
	}
	return outSlice
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
