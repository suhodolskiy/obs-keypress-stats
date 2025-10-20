package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

//go:embed index.html
var tpl embed.FS

type Server struct {
	templatePath string
	srv          *http.Server
	broker       *Broker
}

func NewServer(templatePath string) *Server {
	server := Server{
		broker: NewBroker(),
	}

	if templatePath != "" {
		server.templatePath = templatePath
	}

	return &server
}

func (s *Server) Broadcast(count int) {
	s.broker.Broadcast(count)
}

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ch := s.broker.AddClient()
	defer s.broker.RemoveClient(ch)

	for {
		select {
		case msg := <-ch:
			fmt.Fprintf(w, "data: %d\n\n", msg)
			flusher.Flush()
		case <-s.broker.Done:
			return
		case <-r.Context().Done():
			return
		}
	}
}

func (s *Server) Start(addr string) {
	mux := http.NewServeMux()

	mux.HandleFunc(
		"/", func(w http.ResponseWriter, r *http.Request) {
			if s.templatePath != "" {
				exePath, err := os.Executable()
				if err != nil {
					panic(err)
				}
				exeDir := filepath.Dir(exePath)
				tplPath := filepath.Join(exeDir, s.templatePath)

				http.ServeFile(w, r, tplPath)
			} else {
				data, err := tpl.ReadFile("index.html")
				if err != nil {
					http.Error(w, "File not found", http.StatusNotFound)
					return
				}

				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.Write(data)
			}
		},
	)
	mux.HandleFunc("/events", s.handler)

	s.srv = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Println("Server started at http://localhost:8088/ 🚀")

	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server Listed Failed: %s\n", err)
	}
}

func (s *Server) Stop() {
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	s.broker.Shutdown()

	if err := s.srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	log.Println("Server gracefully stopped 🛑")
}
