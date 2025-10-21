package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	hook "github.com/robotn/gohook"
)

func main() {
	addr := flag.String("addr", "localhost:8088", "Address for the HTTP server to listen on")
	template := flag.String("template", "", "Path to the custom HTML template file (e.g. index.html)")
	initialCount := flag.Int("initial-count", 0, "Initial value of the counter")
	stateFile := flag.String("state-file", "./state.txt", "Path to the state file (used to save and restore state)")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())

	server := NewServer(*template)
	if c := ReadCount(*stateFile); c >= *initialCount {
		*initialCount = c
	}

	go func() {
		s := hook.Start()
		defer hook.End()

		for {
			select {
			case ev := <-s:
				if ev.Kind != hook.KeyDown {
					continue
				}

				*initialCount++
				server.Broadcast(*initialCount)
			case <-ctx.Done():
				log.Println("Key press listener stopped ⌨️")
				return
			}
		}
	}()

	go server.Start(*addr)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	log.Println("Shutting down...")
	cancel()

	PersistCount(*stateFile, *initialCount)
	server.Stop()
}
