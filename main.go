package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/polisgo2020/search-bolotrush/config"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"

	"github.com/polisgo2020/search-bolotrush/index"
	"github.com/polisgo2020/search-bolotrush/web"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		zlog.Err(err).Msg("can not load configs")
	}
	logLvl, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		zlog.Err(err).Msg("can not parse loglvl")
	}
	zerolog.SetGlobalLevel(logLvl)

	fileFlag := flag.Bool("f", false, "save index to file")
	searchFlag := flag.String("s", "", "search query")
	webFlag := flag.Bool("web", false, "input listen interface")
	flag.Parse()
	if flag.NArg() != 1 {
		zlog.Fatal().Err(errors.New("flag error")).Msg("there's wrong number of input arguments")
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
				fmt.Printf("%d) %s: matches - %d\n", i+1, match.Filename, match.Matches)
			}
		} else {
			fmt.Println("There's no results :(")
		}
	}
	if *webFlag {
		server, err := web.NewServer(cfg.Listen, InvertedIndexMap)
		if err != nil {
			zlog.Err(err).Msg("can't create server")
		}
		if err := server.Run(); err != nil {
			zlog.Err(err).Msg("can't run server")
		}
	}
}

func textBuilder(path string, InvertedIndexMap *index.InvMap) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		zlog.Fatal().Err(err).Msg("can not read path")
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
			if err != nil {
				zlog.Fatal().Err(err).Msg("can not read file")
			}

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
	if err != nil {
		zlog.Fatal().Err(err).Msg("can not create out file")
	}

	for key, value := range inputMap {
		strSlice := index.GetDocStrSlice(value)
		_, err := file.WriteString(key + ": {" + strings.Join(strSlice, ",") + "}\n")
		if err != nil {
			zlog.Fatal().Err(err).Msg("can not write text to file")
		}
	}
	err = file.Close()
	if err != nil {
		zlog.Fatal().Err(err).Msg("can not close file")
	}
}
