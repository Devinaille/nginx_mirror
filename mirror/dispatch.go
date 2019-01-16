package mirror

import (
	"errors"
	"fmt"
	"log"
)

// Dispatcher Dispatch workers, max worker number is define by MaxWorkers
type Dispatcher struct {
	MaxWorkers int
	WorkerPool chan chan Payload
	Workers    []*Worker
	quit       chan bool
}

// NewDispatcher new a dispatcher
func NewDispatcher(maxWorkers int) (*Dispatcher, error) {

	if maxWorkers == 0 {
		// no workers, no dispatcher.
		return nil, errors.New("Zero worker numbers")
	}
	return &Dispatcher{
		MaxWorkers: maxWorkers,
		WorkerPool: make(chan chan Payload, maxWorkers),
		Workers:    make([]*Worker, maxWorkers),
		quit:       make(chan bool),
	}, nil
}

//Run run a dispatcher to dispatcher workers
func (d *Dispatcher) Run(queue chan Payload) {
	for i := 0; i < d.MaxWorkers; i++ {
		worker, err := NewWorker(d.WorkerPool)
		if err != nil {
			log.Printf("One worker create failed, %s\n", err.Error())
			continue
		}

		d.Workers[i] = worker
		fmt.Printf("%p %p\n", worker.quit, d.Workers[i].quit)
		worker.Start()
	}
	// fmt.Println(d.QuitPool)
	go d.dispatch(queue)
}

func (d *Dispatcher) dispatch(queue chan Payload) {
	for {
		select {
		case payload := <-queue:
			// fmt.Println("dispatch get a job.")
			// a job request has been received
			go func(p Payload) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				payloadQueue := <-d.WorkerPool

				// dispatch the job to the worker job channel
				payloadQueue <- p
			}(payload)
		case <-d.quit:
			close(queue)
			for _, w := range d.Workers {
				fmt.Printf("%p\n", w.quit)
				w.Stop()
			}
			return
		}

	}
}

//Stop stop the dispatcher
func (d *Dispatcher) Stop() {
	// stop dispatcher first
	log.Println("Stop dispatcher!")
	d.quit <- true
}
