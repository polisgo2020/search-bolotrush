package main

import (
	"errors"
	"fmt"
	"github.com/polisgo2020/search-bolotrush/index"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal(errors.New("not enough program arguments"))
	}
	InvertedIndexMap := index.NewInvMap()

	textBuilder(os.Args[1], &InvertedIndexMap)

	switch os.Args[2] {
	case "index":
		writeMapToFile(&InvertedIndexMap)

	case "search":
		if len(os.Args) < 4 {
			fmt.Println("There's nothing to search")
			return
		}
		matchListOut := index.Searcher(&InvertedIndexMap, os.Args[3:])
		fmt.Println("Search result:")
		if len(matchListOut) > 0 {
			for i, match := range matchListOut {
				if i > 4 {
					break
				}
				fmt.Printf("%d) %s: matches - %d\n", i+1, match.FileName, match.Matches)
			}
		} else {
			fmt.Println("There's no results :(")
		}

	default:
		fmt.Println("Command is unknown. Try again.")
		return
	}
}

func textBuilder(path string, InvertedIndexMap *index.InvMap) {

	files, err := ioutil.ReadDir(path)
	checkError(err)

	textChannel := make(chan index.StraightIndex)
	wg := &sync.WaitGroup{}
	mutex := &sync.Mutex{}

	go index.AsyncInvertIndex(textChannel, InvertedIndexMap, mutex, wg)

	regCompiled := regexp.MustCompile(`[\W]+`)
	for _, file := range files {
		wg.Add(1)
		go asyncRead(file, regCompiled, textChannel, path, wg)
	}
	wg.Wait()
	close(textChannel)
}

func asyncRead(file os.FileInfo, reg *regexp.Regexp, ch chan<- index.StraightIndex, path string, wg *sync.WaitGroup) {
	defer wg.Done()
	text, err := ioutil.ReadFile(path + "/" + file.Name())
	checkError(err)
	chStruct := index.StraightIndex{
		FileName: strings.TrimRight(file.Name(), ".txt"),
		Text:     reg.Split(string(text), -1),
	}
	ch <- chStruct
	fmt.Println("Вызвал " + file.Name())
}

func writeMapToFile(inputMap *index.InvMap) {

	file, err := os.Create("out.txt")
	checkError(err)

	for key, value := range *inputMap {
		strSlice := index.GetDocStrSlice(value)
		_, err := file.WriteString(key + ": {" + strings.Join(strSlice, ",") + "}\n")
		checkError(err)
	}
	err = file.Close()
	checkError(err)
}

func checkError(err error) {

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
