package util

import (
	"github.com/mustafaturan/bus"
	log "github.com/sirupsen/logrus"
)

const (
	DefaultNumWorkers = 4
	DefaultBufferSize = 10000
)

type Pool struct {
	workers []*Worker
	events  chan *bus.Event
	quit    chan bool
}

func NewPool(numWorkers int, bufferSize int) *Pool {
	events := make(chan *bus.Event, bufferSize)
	quit := make(chan bool, numWorkers)
	workers := make([]*Worker, numWorkers)
	for index := 0; index < numWorkers; index++ {
		workers[index] = &Worker{events: events}
	}

	return &Pool{
		workers: workers,
		events:  events,
		quit:    quit,
	}
}

func (p *Pool) SetFunc(f func(e *bus.Event)) {
	for _, w := range p.workers {
		go w.Start(f)
	}
}

func (p *Pool) Consume(e *bus.Event) {
	p.events <- e
}

func (p *Pool) Close() {
	for range p.workers {
		p.quit <- true
	}
}

type Worker struct {
	events chan *bus.Event
	quit   chan bool
}

func (w *Worker) Start(f func(e *bus.Event)) {
	for {
		select {
		case e := <-w.events:
			f(e)
		case <-w.quit:
			log.Info("shutdown worker")
			break
		}
	}
}
