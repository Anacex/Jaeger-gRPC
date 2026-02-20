package main

import (
	"context"
	"log"
	"sync"
)

type WorkerPool struct {
	queue   chan LogEntry
	handler LogHandler
	wg      sync.WaitGroup
}

func NewWorkerPool(handler LogHandler, bufferSize int, workers int) *WorkerPool {
	wp := &WorkerPool{
		queue:   make(chan LogEntry, bufferSize),
		handler: handler,
	}

	for i := 0; i < workers; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}
	return wp
}

func (wp *WorkerPool) worker() {
	defer wp.wg.Done()

	for entry := range wp.queue {

		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Println("Recovered from panic in worker:", r)
				}
			}()

			err := wp.handler.Handle(context.Background(), entry)
			if err != nil {
				log.Println("Failed to process log:", err)
			}
		}()
	}
}

func (wp *WorkerPool) Enqueue(entry LogEntry) bool {
	select {
	case wp.queue <- entry:
		return true
	default:
		log.Println("Queue full, dropping log:", entry)
		return false
	}
}
