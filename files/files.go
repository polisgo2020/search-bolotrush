package files

import (
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/polisgo2020/search-bolotrush/index"

	zl "github.com/rs/zerolog/log"
)

func IndexBuilder(path string /*InvertedIndexMap *index.InvMap*/) (index.InvMap, error) {
	indexMap := &index.InvMap{}

	channel := make(chan index.StraightIndex)
	wg := &sync.WaitGroup{}
	mutex := &sync.Mutex{}
	go func(channel chan index.StraightIndex, indexMap *index.InvMap) {
		for input := range channel {
			input.Wg.Add(1)
			input.Mutex.Lock()
			indexMap.InvertIndex(input.Text, input.FileName)
			input.Mutex.Unlock()
			input.Wg.Done()
		}
	}(channel, indexMap)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	wg.Add(len(files))
	for _, file := range files {
		go asyncRead(wg, path, file, mutex, channel)
	}
	wg.Wait()
	close(channel)
	return *indexMap, nil
}

func asyncRead(wg *sync.WaitGroup, path string, file os.FileInfo, mutex *sync.Mutex, channel chan index.StraightIndex) {
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
}

func WriteMapToFile(inputMap index.InvMap) error {
	file, err := os.Create("out.txt")
	if err != nil {
		return err
	}

	for key, value := range inputMap {
		strSlice := index.GetDocStrSlice(value)
		_, err := file.WriteString(key + ": {" + strings.Join(strSlice, ",") + "}\n")
		if err != nil {
			return err
		}

	}
	err = file.Close()
	if err != nil {
		return err
	}
	return nil
}

func checkError(err error) {
	if err != nil {
		zl.Fatal().Err(err)
	}
}
