package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"
)

const minWordLength = 2

var (
	stopWORDS = map[string]bool{
		"":    true,
		" ":   true,
		"a":   true,
		"an":  true,
		"the": true,
	}
	invertedIndexMap = make(map[string][]string)
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("File path is not in program arguments!")
		return
	}
	fileToText(os.Args[1])
	writeMapToFile(invertedIndexMap)
}

func fileToText(path string) {

	files, err := ioutil.ReadDir(path)
	checkError(err)
	regCompiled := regexp.MustCompile(`[\W]+`)

	for _, file := range files {

		text, err := ioutil.ReadFile(path + "/" + file.Name())
		checkError(err)

		regularText := regCompiled.Split(string(text), -1)
		invertIndex(regularText, strings.TrimRight(file.Name(), ".txt"))
	}
}

func writeMapToFile(inputMap map[string][]string) {

	file, err := os.Create("out.txt")
	checkError(err)
	defer file.Close()

	for key, value := range inputMap {
		_, err := file.WriteString(key + ": " + "{" + strings.Join(value, ",") + "}\n")
		checkError(err)
	}
}

func invertIndex(inputWords []string, docId string) {

	inputWords = cleanText(inputWords)
	for _, word := range inputWords {
		if !isStringInSlice(docId, invertedIndexMap[word]) {
			invertedIndexMap[word] = append(invertedIndexMap[word], docId)
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

func checkError(err error) {

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
