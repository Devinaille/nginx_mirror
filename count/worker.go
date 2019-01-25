package count

import (
	"log"
	"nginx_mirror/mirror"

	"github.com/google/uuid"
)

// Worker 用于计数工作
type Worker struct {
	WorkerID   uuid.UUID
	WorkerPool chan chan mirror.Request
	Input      chan mirror.Request
	Counter    Counter
	status     bool
	quit       chan bool
}

// NewWorker 初始化计数worker
func NewWorker(workerPool chan chan mirror.Request, c Counter) (*Worker, error) {
	workerID, err := uuid.NewRandom()
	if err != nil {
		log.Printf("Get uuid failed, err: %s\n", err.Error())
		workerID = uuid.New()
	}
	return &Worker{
		WorkerID:   workerID,
		WorkerPool: workerPool,
		Counter:    c,
		Input:      make(chan mirror.Request),
		status:     false,
		quit:       make(chan bool),
	}, nil
}

// Start 处理mirror发来的Request内容。
func (w *Worker) Start() {
	log.Printf("Count worker %s start!\n", w.WorkerID)
	w.status = true
	go func() {
		for {
			w.WorkerPool <- w.Input
			// log.Printf("Count worker %s wait for payload.\n", w.WorkerID)
			select {
			case request := <-w.Input:
				w.Counter.Add(request.XOriginalURI)
			case <-w.quit:
				log.Printf("Count Worker %s quit!\n", w.WorkerID)
				return
			}
		}
	}()
}

//Stop stop a countWorker
func (w *Worker) Stop() {
	go func() {
		log.Printf("Stop worker %s!\n", w.WorkerID)
		w.status = false
		w.quit <- true
	}()
}

//ID 返回workerID
func (w *Worker) ID() string {
	return w.WorkerID.String()
}

//Status 查询Worker的状态
func (w *Worker) Status() bool {
	return w.status
}
