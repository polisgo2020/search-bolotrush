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
	webFlag := flag.String("web", "", "input port")
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatal(errors.New("there's wrong number of input arguments"))
	}

	InvertedIndexMap := index.NewInvMap()
	if err := textBuilder(flag.Args()[0], &InvertedIndexMap); err != nil {
		log.Fatal(errors.New("path to files is incorrect"))
	}
	if *fileFlag {
		writeMapToFile(InvertedIndexMap)
	}
	if *searchFlag != "" {
		matchListOut := InvertedIndexMap.Searcher(strings.Fields(*searchFlag))
		index.ShowSearchResults(matchListOut)
	}
	if *webFlag != "" {
		if err := web.RunServer(":"+*webFlag, InvertedIndexMap); err != nil {
			log.Fatal(err)
		}
	}
}

func textBuilder(path string, InvertedIndexMap *index.InvMap) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
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
				FileName: strings.TrimRight(file.Name(), ".txt"),
				Text:     string(text),
				Wg:       wg,
				Mutex:    mutex,
			}
			channel <- info
		}(file)
	}
	wg.Wait()
	close(channel)
	return nil
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
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
