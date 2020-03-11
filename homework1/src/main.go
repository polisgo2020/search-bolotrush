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
	stopWORDS = []string{
		"",
		" ",
		"a",
		"an",
		"the",
	}
	invertedIndexMap = make(map[string][]string)
)

func main() {
	setFileToText("homework1/books")
	/*	for key, _ := range invertedIndexMap {
			fmt.Println(key + ": " + strconv.Itoa(utf8.RuneCountInString(key)))
		}
	*/
	writeMapToFile(invertedIndexMap)
}

func setFileToText(path string) {

	files, err := ioutil.ReadDir(path)
	checkError(err)

	for _, file := range files {

		text, err := ioutil.ReadFile(path + "/" + file.Name())
		checkError(err)

		regularText := regexp.MustCompile(`[\W]+`).Split(string(text), -1)
		regularText = cleanText(regularText)
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

func invertIndex(inputWords []string, docId string) /*[]string*/ {
	for _, word := range inputWords {
		if !stringInSlice(docId, invertedIndexMap[word]) {
			invertedIndexMap[word] = append(invertedIndexMap[word], docId)
		}
	}
}

func cleanText(inputWords []string) []string {

	cleanWords := make([]string, 0)
	for _, word := range inputWords {

		if stringInSlice(word, stopWORDS) || utf8.RuneCountInString(word) < minWordLength {
			continue
		}
		word = strings.ToLower(word)
		cleanWords = append(cleanWords, word)
	}
	return cleanWords
}

func stringInSlice(a string, list []string) bool {
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
