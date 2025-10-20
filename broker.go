package main

import (
	"sync"
)

type Broker struct {
	clients map[chan int]struct{}
	lock    sync.Mutex

	Done chan struct{}
}

func NewBroker() *Broker {
	return &Broker{
		clients: make(map[chan int]struct{}),
		Done:    make(chan struct{}),
	}
}

func (b *Broker) AddClient() chan int {
	b.lock.Lock()
	defer b.lock.Unlock()
	ch := make(chan int)
	b.clients[ch] = struct{}{}
	return ch
}

func (b *Broker) RemoveClient(ch chan int) {
	b.lock.Lock()
	defer b.lock.Unlock()
	delete(b.clients, ch)
}

func (b *Broker) Broadcast(count int) {
	b.lock.Lock()
	defer b.lock.Unlock()

	for ch := range b.clients {
		select {
		case ch <- count:
		}
	}
}

func (b *Broker) Shutdown() {
	close(b.Done)
}
