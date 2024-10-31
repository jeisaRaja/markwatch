package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	ws = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	connections []*websocket.Conn
	mu          sync.Mutex
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

func (s *Server) Start() error {
	http.HandleFunc("/", s.fileHandler)
	http.HandleFunc("/ws", s.handleWs)
	return s.s.ListenAndServe()
}

func (s *Server) fileHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, s.fp)
}

func (s *Server) reload(fname string) error {
	if _, err := os.Stat(fname); err != nil {
		return ErrFileNotFound
	}

	s.fp = fname

	return nil
}

func (s *Server) handleWs(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer conn.Close()

	mu.Lock()
	connections = append(connections, conn)
	mu.Unlock()

	defer func() {
		mu.Lock()
		for i, c := range connections {
			if c == conn {
				connections = append(connections[:i], connections[i+1:]...)
				break
			}
		}
		mu.Unlock()
		conn.Close()
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func triggerRefresh() {
	mu.Lock()
	for _, conn := range connections {
		err := conn.WriteMessage(websocket.TextMessage, []byte("refresh"))
		if err != nil {
			conn.Close()
		}
	}
	mu.Unlock()
}
