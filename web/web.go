package web

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/polisgo2020/search-bolotrush/index"
	zlog "github.com/rs/zerolog/log"
)

type Server struct {
	server         http.Server
	index          index.InvMap
	startTemplate  *template.Template
	searchTemplate *template.Template
}

func NewServer(addr string, index index.InvMap) (*Server, error) {
	if index == nil {
		return nil, errors.New("incorrect index")
	}

	if addr == "" {
		return nil, errors.New("incorrect address")
	}
	startHTML, err := template.ParseFiles("web/start.html")
	if err != nil {
		zlog.Fatal().Err(err).Msg("can not read index template")
	}
	searchHTML, err := template.ParseFiles("web/search.html")
	if err != nil {
		zlog.Fatal().Err(err).Msg("can not read index template")
	}
	s := &Server{
		index:          index,
		startTemplate:  startHTML,
		searchTemplate: searchHTML,
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
	zlog.Debug().Msgf("starting server at %s", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *Server) searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		_, err := fmt.Fprintln(w, "Wrong query")
		if err != nil {
			zlog.Fatal().Err(err)
		}
	}
	result := s.index.Searcher(strings.Fields(query))
	if len(result) == 0 {
		_, err := fmt.Fprintln(w, "There's no results :(")
		if err != nil {
			zlog.Fatal().Err(err)
		}
		return
	}
	err := s.searchTemplate.Execute(w, struct {
		Result []index.MatchList
		Query  string
	}{
		Result: result,
		Query:  query,
	})
	if err != nil {
		zlog.Fatal().Err(err).Msg("can not render template")
	}
}

func (s *Server) startHandler(w http.ResponseWriter, r *http.Request) {
	err := s.startTemplate.Execute(w, nil)
	if err != nil {
		zlog.Fatal().Err(err).Msg("can not render template")
	}
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		zlog.Debug().
			Str("method", r.Method).
			Str("remote", r.RemoteAddr).
			Str("path", r.URL.Path).
			Int("duration", int(time.Since(start))).
			Msgf("called url %s", r.URL.Path)
	})
}
