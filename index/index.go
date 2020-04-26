package index

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
	"unicode/utf8"
)

const minWordLength = 2

var regCompiled = regexp.MustCompile(`[^a-zA-Z_]+`)

type WordInfo struct {
	Filename  string
	Positions []int
}

type InvMap map[string][]WordInfo

type StraightIndex struct {
	Filename string
	Text     string
	Mutex    *sync.Mutex
	Wg       *sync.WaitGroup
}

func NewInvMap() InvMap {
	index := make(InvMap)
	return index
}

func (thisMap *InvMap) AsyncInvertIndex(docChan chan StraightIndex) {
	for input := range docChan {
		input.Wg.Add(1)
		input.Mutex.Lock()
		thisMap.InvertIndex(input.Text, input.Filename)
		input.Mutex.Unlock()
		input.Wg.Done()
	}
}

func (thisMap *InvMap) InvertIndex(inputText string, fileName string) {
	wordList := prepareText(inputText)
	for i, word := range wordList {
		if index, ok := thisMap.isWordInList(word, fileName); !ok {
			structure := WordInfo{
				Filename:  fileName,
				Positions: []int{},
			}
			structure.Positions = append(structure.Positions, i)
			(*thisMap)[word] = append((*thisMap)[word], structure)
		} else if index != -1 {
			(*thisMap)[word][index].Positions = append((*thisMap)[word][index].Positions, i)
		}

	}
}

func GetDocStrSlice(slice []WordInfo) []string {
	outSlice := make([]string, 0)
	for _, doc := range slice {
		outSlice = append(outSlice, doc.Filename)
	}
	return outSlice
}

type MatchList struct {
	Matches  int
	Filename string
}

func (thisMap InvMap) Searcher(query []string) []MatchList {
	var matchesSlice []MatchList
	var matchesMap = make(map[string]int, 0)
	query = cleanText(query)
	for _, word := range query {
		if fileList, ok := thisMap[word]; ok {
			for _, fileName := range fileList {
				matchesMap[fileName.Filename] += len(fileName.Positions)
			}
		}
	}
	for name, matches := range matchesMap {
		matchesSlice = append(matchesSlice, MatchList{
			Matches:  matches,
			Filename: name,
		})
	}
	if len(matchesSlice) > 0 {
		sort.Slice(matchesSlice, func(i, j int) bool {
			return matchesSlice[i].Matches > matchesSlice[j].Matches
		})
	}
	return matchesSlice
}

func ShowSearchResults(matchListOut []MatchList) {
	fmt.Println("Search result:")
	if len(matchListOut) > 0 {
		for i, match := range matchListOut {
			if i > 4 {
				break
			}
			fmt.Printf("%d) %s: matches - %d\n", i+1, match.Filename, match.Matches)
		}
	} else {
		fmt.Println("There's no results :(")
	}
}

func (thisMap InvMap) isWordInList(word string, docId string) (int, bool) {
	for i, ind := range thisMap[word] {
		if ind.Filename == docId {
			return i, true
		}
	}
	return -1, false
}

func prepareText(in string) []string {
	tokens := cleanText(regCompiled.Split(in, -1))
	return tokens
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
