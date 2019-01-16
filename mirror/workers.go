package mirror

import (
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
)

// Payload mirror转发过来的载荷
type Payload struct {
	Headers http.Header
	Method  string
}

// Worker 处理mirror转发过来的请求
type Worker struct {
	WorkerID   uuid.UUID
	WorkerPool chan chan Payload
	Queue      chan Payload
	// quitPool   chan chan bool
	quit chan bool
}

// NewWorker new a worker
func NewWorker(workerPool chan chan Payload) (*Worker, error) {
	workerID, err := uuid.NewRandom()
	if err != nil {
		log.Printf("Get uuid failed, err: %s\n", err.Error())
		return nil, errors.New("Get worker id failed")
	}
	return &Worker{
		WorkerID:   workerID,
		WorkerPool: workerPool,
		Queue:      make(chan Payload),
		quit:       make(chan bool),
	}, nil
}

// Start start a worker for mirror
func (w *Worker) Start() {
	log.Printf("Worker %s start!\n", w.WorkerID)
	go func() {
		for {
			log.Printf("Worker %s start wait for payload.\n", w.WorkerID)
			w.WorkerPool <- w.Queue
			log.Printf("Worker %s wait for payload.\n", w.WorkerID)
			select {
			case payload := <-w.Queue:
				for header, value := range payload.Headers {
					log.Printf("Header: %s, value %+v\n", header, value)
				}
				// log.Printf("%s get payload. %+v!\n", w.WorkerID, payload)
			case <-w.quit:
				log.Printf("Worker %s quit!\n", w.WorkerID)
				return

			}
		}
	}()
}

//Stop stop a worker
func (w *Worker) Stop() {
	go func() {
		log.Printf("Stop worker %s!\n", w.WorkerID)
		w.quit <- true
	}()
}
