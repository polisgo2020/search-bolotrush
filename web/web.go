package web

import (
	"errors"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/polisgo2020/search-bolotrush/index"
	zl "github.com/rs/zerolog/log"
)

type Server struct {
	server         http.Server
	index          index.InvMap
	startTemplate  *template.Template
	searchTemplate *template.Template
	searchFunc     func(query string) ([]index.MatchList, error)
}

func NewServer(addr string, searcher func(query string) ([]index.MatchList, error)) (*Server, error) {
	if addr == "" {
		return nil, errors.New("incorrect address")
	}
	startHTML, err := template.ParseFiles("web/start.html")
	if err != nil {
		return nil, fmt.Errorf("can't read index template")
	}
	searchHTML, err := template.ParseFiles("web/search.html")
	if err != nil {
		return nil, fmt.Errorf("can't read index template")
	}
	s := &Server{
		startTemplate:  startHTML,
		searchTemplate: searchHTML,
		searchFunc:     searcher,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.startHandler)
	mux.HandleFunc("/search", s.searchHandler)

	logServer := logger(mux)
	s.server = http.Server{
		Addr:         addr,
		Handler:      logServer,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return s, nil
}

func (s *Server) Run() error {
	zl.Debug().Msgf("starting server at %s", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *Server) searchHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query().Get("query")
	if query == "" {
		fmt.Fprintln(w, "Wrong query")
	}
	result, err := s.searchFunc(query)
	if err != nil {
		zl.Error().Err(err).Msg("error while searching")
		return
	}
	if len(result) == 0 {
		fmt.Fprintln(w, "There's no results :(")
	}
	err = s.searchTemplate.Execute(w, struct {
		Result []index.MatchList
		Query  string
	}{
		Result: result,
		Query:  query,
	})
	if err != nil {
		zl.Error().Err(err).Msg("can not render template")
	}
}

func (s *Server) startHandler(w http.ResponseWriter, r *http.Request) {
	err := s.startTemplate.Execute(w, nil)
	if err != nil {
		zl.Error().Err(err).Msg("can not render template")
	}
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		zl.Info().
			Str("method", r.Method).
			Str("remote", r.RemoteAddr).
			Str("path", r.URL.Path).
			Int("duration", int(time.Since(start))).
			Msgf("called url %s", r.URL.Path)
	})
}
