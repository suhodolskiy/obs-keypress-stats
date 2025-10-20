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
	addr := flag.String("addr", ":8088", "HTTP server address")
	tplPath := flag.String("tpl", "", "Template path")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())

	var (
		count  = ReadCount("./count.txt")
		server = NewServer(*tplPath)
	)

	go func() {
		s := hook.Start()
		defer hook.End()

		for {
			select {
			case ev := <-s:
				if ev.Kind != hook.KeyDown {
					continue
				}

				count++
				server.Broadcast(count)
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

	PersistCount("./count.txt", count)
	server.Stop()
}
