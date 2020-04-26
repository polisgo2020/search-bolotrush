package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/polisgo2020/search-bolotrush/index"
	"github.com/polisgo2020/search-bolotrush/web"
)

func main() {
	fileFlag := flag.Bool("f", false, "save index to file")
	searchFlag := flag.String("s", "", "search query")
	webFlag := flag.String("web", "", "input listen interface")
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatal(errors.New("there's wrong number of input arguments"))
	}

	InvertedIndexMap := index.NewInvMap()
	textBuilder(flag.Args()[0], &InvertedIndexMap)
	if *fileFlag {
		writeMapToFile(InvertedIndexMap)
	}
	if *searchFlag != "" {
		matchListOut := InvertedIndexMap.Searcher(strings.Fields(*searchFlag))
		fmt.Println("Search result:")
		if len(matchListOut) > 0 {
			for i, match := range matchListOut {
				fmt.Printf("%d) %s: matches - %d\n", i+1, match.FileName, match.Matches)
			}
		} else {
			fmt.Println("There's no results :(")
		}
	}
	if *webFlag != "" {
		server, err := web.NewServer(*webFlag, InvertedIndexMap)
		if err != nil {
			log.Fatal(err)
		}
		if err := server.Run(); err != nil {
			log.Fatal(err)
		}
	}
}

func textBuilder(path string, InvertedIndexMap *index.InvMap) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	channel := make(chan index.StraightIndex)
	wg := &sync.WaitGroup{}
	mutex := &sync.Mutex{}
	go InvertedIndexMap.AsyncInvertIndex(channel)

	for _, file := range files {
		wg.Add(1)
		go func(file os.FileInfo) {
			defer wg.Done()
			text, err := ioutil.ReadFile(path + "/" + file.Name())
			checkError(err)

			info := index.StraightIndex{
				Filename: strings.TrimRight(file.Name(), ".txt"),
				Text:     string(text),
				Wg:       wg,
				Mutex:    mutex,
			}
			channel <- info
		}(file)
	}
	wg.Wait()
	close(channel)
}

func writeMapToFile(inputMap index.InvMap) {
	file, err := os.Create("out.txt")
	checkError(err)

	for key, value := range inputMap {
		strSlice := index.GetDocStrSlice(value)
		_, err := file.WriteString(key + ": {" + strings.Join(strSlice, ",") + "}\n")
		checkError(err)
	}
	err = file.Close()
	checkError(err)
}

func checkError(err error) {

	if err != nil {
		log.Fatal(err.Error())
	}
}
