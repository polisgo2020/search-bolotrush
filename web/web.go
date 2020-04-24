package web

import (
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/polisgo2020/search-bolotrush/index"
	zl "github.com/rs/zerolog/log"
)

func RunServer(index index.InvMap, addr string) {
	startHTML, err := template.ParseFiles("web/start.html")
	if err != nil {
		zl.Fatal().Err(err).Msg("can not read index template")
	}
	searchHTML, err := template.ParseFiles("web/search.html")
	if err != nil {
		zl.Fatal().Err(err).Msg("can not read index template")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		startHandler(w, startHTML)
	})
	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		searchHandler(w, r, index, searchHTML)
	})
	logServer := logger(mux)
	server := http.Server{
		Addr:         addr,
		Handler:      logServer,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	zl.Debug().Msgf("starting server at", addr)
	err = server.ListenAndServe()
	if err != nil {
		zl.Fatal().Err(err)
	}

}

func searchHandler(w http.ResponseWriter, r *http.Request, inputIndex index.InvMap, template *template.Template) {
	query := r.URL.Query().Get("query")
	if query == "" {
		_, err := fmt.Fprintln(w, "Wrong query")
		if err != nil {
			zl.Fatal().Err(err)
		}
		return
	}
	result := inputIndex.Search(query)
	if len(result) == 0 {
		_, err := fmt.Fprintln(w, "There's no results :(")
		if err != nil {
			zl.Fatal().Err(err)
		}
		return
	}
	err := template.Execute(w, struct {
		Result []index.MatchList
		Query  string
	}{
		Result: result,
		Query:  query,
	})
	if err != nil {
		zl.Fatal().Err(err).Msg("can not render template")
	}
}

func startHandler(w http.ResponseWriter, template *template.Template) {
	err := template.Execute(w, nil)
	if err != nil {
		zl.Fatal().Err(err).Msg("can not render template")
	}
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		zl.Debug().
			Str("method", r.Method).
			Str("remote", r.RemoteAddr).
			Str("path", r.URL.Path).
			Int("duration", int(time.Since(start))).
			Msgf("called url %s", r.URL.Path)
	})
}
