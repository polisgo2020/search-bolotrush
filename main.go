package main

import (
	"errors"
	"fmt"
	"os"
	_ "strings"

	"github.com/polisgo2020/search-bolotrush/web"

	"github.com/polisgo2020/search-bolotrush/files"

	"github.com/urfave/cli/v2"

	"github.com/rs/zerolog"

	zl "github.com/rs/zerolog/log"

	"github.com/polisgo2020/search-bolotrush/config"
)

var cfg config.Config

func main() {
	var err error
	cfg, err = config.Load()
	if err != nil {
		zl.Fatal().Err(err)
	}

	loglevel, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		zl.Fatal().Err(err)
	}

	zerolog.SetGlobalLevel(loglevel)

	app := &cli.App{
		Name:  "Searcher",
		Usage: "The app searches docs using inverted index and find the best match",
	}
	pathFlag := &cli.StringFlag{
		Name:     "path",
		Aliases:  []string{"p"},
		Usage:    "Path to files directory",
		Required: true,
	}
	app.Commands = []*cli.Command{
		{
			Name:    "file",
			Aliases: []string{"f"},
			Usage:   "Saves index to file",
			Flags:   []cli.Flag{pathFlag},
			Action:  indexFile,
		}, {
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "Reads query and search files",
			Flags:   []cli.Flag{pathFlag},
			Subcommands: []*cli.Command{
				{
					Name:    "console",
					Aliases: []string{"c"},
					Usage:   "Searches index in console",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "query",
							Aliases:  []string{"q"},
							Usage:    "Searches query in console",
							Required: true,
						},
					},
					Action: searchConsole,
				}, {
					Name:    "web",
					Aliases: []string{"w"},
					Usage:   "Creates web server for search using http",
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:     "db",
							Usage:    "Uses PostgreSQL for data",
							Required: true,
						},
						&cli.StringFlag{
							Name:     "port",
							Usage:    "Web server's port",
							Required: true,
						},
					},
					Action: searchWeb,
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		zl.Fatal().Err(err)
	}
}

func indexFile(c *cli.Context) error {
	path := c.String("path")
	if len(path) == 0 {
		return errors.New("path to files is not found")
	}

	indexMap, err := files.IndexBuilder(path)
	if err != nil {
		return err
	}

	err = files.WriteMapToFile(indexMap)
	if err != nil {
		return err
	}
	return nil
}

func searchConsole(c *cli.Context) error {
	path := c.String("path")
	if len(path) == 0 {
		return errors.New("path to files is not found")
	}

	query := c.String("query")
	if len(path) == 0 {
		return errors.New("query phrase is not found")
	}

	indexMap, err := files.IndexBuilder(path)
	if err != nil {
		return err
	}

	matches := indexMap.Search(query)
	if len(matches) > 0 {
		for i, match := range matches {
			fmt.Printf("%d) %s: matches - %d\n", i+1, match.FileName, match.Matches)
		}
	} else {
		fmt.Println("There's no results :(")
	}
	return nil
}

func searchWeb(c *cli.Context) error {
	path := c.String("path")
	if len(path) == 0 {
		return errors.New("path to files is not found")
	}

	indexMap, err := files.IndexBuilder(path)
	if err != nil {
		return err
	}

	port := c.String("port")
	web.RunServer(indexMap, port)
	return nil
}
