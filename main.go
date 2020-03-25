package main

import (
	"fmt"
	"github.com/polisgo2020/search-bolotrush/index"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Not enough program arguments!")
		return
	}
	fileToText(os.Args[1])

	switch os.Args[2] {
	case "index":

		writeMapToFile(index.InvertedIndexMap)

	case "search":

		if len(os.Args) < 4 {
			fmt.Println("There's nothing to search")
			return
		}
		matchListOut := index.Searcher(index.InvertedIndexMap, os.Args[3:])
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

func fileToText(path string) {

	files, err := ioutil.ReadDir(path)
	checkError(err)

	regCompiled := regexp.MustCompile(`[\W]+`)
	for _, file := range files {

		text, err := ioutil.ReadFile(path + "/" + file.Name())
		checkError(err)

		regularText := regCompiled.Split(string(text), -1)

		index.InvertIndex(regularText, strings.TrimRight(file.Name(), ".txt"))
	}
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
