package search

import (
	"github.com/polisgo2020/search-bolotrush/index"
	"sort"
)

type MatchList struct {
	Matches  int
	FileName string
}

func Searcher(invertMap index.InvMap, arguments []string) []MatchList {

	matchesSlice := make([]MatchList, 0)
	matchesMap := make(map[string]int)
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
