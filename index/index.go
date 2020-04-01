package index

import (
	"fmt"
	"sort"
	"strings"
	"sync"
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

type StraightIndex struct {
	FileName string
	Text     []string
}

func (m InvMap) isWordInList(word string, docId string) (int, bool) {
	for i, ind := range m[word] {
		if ind.Doc == docId {
			return i, true
		}
	}
	return -1, false
}

func NewInvMap() InvMap {
	index := make(InvMap)
	return index
}

func AsyncInvertIndex(docChan chan StraightIndex, myMap *InvMap, mutex *sync.Mutex, wg *sync.WaitGroup) {

	for input := range docChan {
		wg.Add(1)
		inputWords := input.Text
		docId := input.FileName
		inputWords = cleanText(inputWords)

		fmt.Println("document: " + docId)

		for i, word := range inputWords {
			mutex.Lock()

			if index, ok := (*myMap).isWordInList(word, docId); !ok {

				structure := wordStruct{
					Doc:      docId,
					Position: []int{},
				}

				structure.Position = append(structure.Position, i)
				(*myMap)[word] = append((*myMap)[word], structure)
			} else if index != -1 {
				(*myMap)[word][index].Position = append((*myMap)[word][index].Position, i)
			}
			mutex.Unlock()
		}
		wg.Done()
	}
}

func GetDocStrSlice(slice []wordStruct) []string {
	outSlice := make([]string, 0)
	for _, doc := range slice {
		outSlice = append(outSlice, doc.Doc)
	}
	return outSlice
}

func Searcher(invertMap *InvMap, arguments []string) []MatchList {

	var matchesSlice []MatchList
	matchesMap := make(map[string]int, 0)

	arguments = cleanText(arguments)
	for _, word := range arguments {
		if docNames, ok := (*invertMap)[word]; ok {
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
