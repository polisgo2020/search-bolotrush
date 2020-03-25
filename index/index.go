package index

import (
	"sort"
	"strings"
	"unicode/utf8"
)

const minWordLength = 2

type wordStruct struct {
	Doc      string
	Position []int
}

type InvMap map[string][]wordStruct

type MatchList struct {
	Matches  int
	FileName string
}

func (p InvMap) isWordInList(word string, docId string) (int, bool) {
	for i, ind := range p[word] {
		if ind.Doc == docId {
			return i, true
		}
	}
	return -1, false
}

var InvertedIndexMap = make(InvMap)

func InvertIndex(inputWords []string, docId string) {

	inputWords = cleanText(inputWords)
	for i, word := range inputWords {
		if index, ok := InvertedIndexMap.isWordInList(word, docId); !ok {

			structure := wordStruct{
				Doc:      docId,
				Position: []int{},
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

func Searcher(invertMap InvMap, arguments []string) []MatchList {

	var matchesSlice []MatchList
	matchesMap := make(map[string]int, 0)

	arguments = cleanText(arguments)
	for _, word := range arguments {
		if docNames, ok := invertMap[word]; ok {
			for _, doc := range docNames {
				matchesMap[doc.Doc] += len(doc.Position)
			}
		}
	}
	for name, matches := range matchesMap {
		matchesSlice = append(matchesSlice, MatchList{
			Matches:  matches,
			FileName: name,
		})
	}
	if len(matchesSlice) > 0 {
		sort.Slice(matchesSlice, func(i, j int) bool {
			return matchesSlice[i].Matches > matchesSlice[j].Matches
		})
	}
	return matchesSlice
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
