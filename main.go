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

	if len(os.Args) < 2 {
		fmt.Println("File path is not in program arguments!")
		return
	}
	fileToText(os.Args[1])
	writeMapToFile(index.InvertedIndexMap)
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

func writeMapToFile(inputMap map[string][]string) {

	file, err := os.Create("out.txt")
	checkError(err)
	defer file.Close()

	for key, value := range inputMap {
		_, err := file.WriteString(key + ": " + "{" + strings.Join(value, ",") + "}\n")
		checkError(err)
	}
}

func checkError(err error) {

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
