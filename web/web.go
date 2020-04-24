package web

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/polisgo2020/search-bolotrush/index"
)

func webSearch(w http.ResponseWriter, r *http.Request, index index.InvMap) {
	query := r.URL.Query().Get("query")
	if query == "" {
		_, err := fmt.Fprintln(w, "Search phrase is not found")
		if err != nil {
			return
		}
		return
	}

	_, err := fmt.Fprintln(w, "Search query: ", query)
	if err != nil {
		return
	}
	result := index.Searcher(strings.Fields(query))
	if len(result) > 0 {
		for i, match := range result {
			_, err := fmt.Fprintf(w, "%d) %s: matches - %d\n", i+1, match.FileName, match.Matches)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		_, err := fmt.Fprintln(w, "There's no results :(")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func RunServer(addr string, index index.InvMap) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		webSearch(w, r, index)
	})
	server := http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Println("starting server at", addr)
	return server.ListenAndServe()
}
