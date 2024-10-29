package main

import (
	"fmt"
	"net/http"
	"os"
)

type Server struct {
	s  *http.Server
	fp string
}

func NewServer() *Server {
	s := &http.Server{
		Addr: ":9090",
	}

	return &Server{s: s}
}

var ErrFileNotFound = fmt.Errorf("File not found")

func (s *Server) reload(fname string) error {
	if _, err := os.Stat(fname); err != nil {
		return ErrFileNotFound
	}

	s.fp = fname

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, s.fp)
	})

	return nil
}
