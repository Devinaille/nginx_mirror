package mirror

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Payload mirror转发过来的载荷
type Payload struct {
	Time    time.Time
	Headers http.Header
	Method  string
}

// Worker 处理mirror转发过来的请求
type Worker struct {
	WorkerID   uuid.UUID
	WorkerPool chan chan Payload
	Input      chan Payload
	Outputs    []chan Request
	quit       chan bool
}

// Request 发送到镜像的请求
type Request struct {
	Time            time.Time
	XRealIP         string
	XOriginalURI    string
	XForwardFor     []string
	XForwarderProto string
	NginxIP         string
	UserAgent       []string
	Method          string
}

// NewWorker new a worker
func NewWorker(workerPool chan chan Payload, o []chan Request) (*Worker, error) {
	workerID, err := uuid.NewRandom()
	if err != nil {
		log.Printf("Get uuid failed, err: %s\n", err.Error())
		return nil, errors.New("Get mirrorWorker id failed")
	}
	return &Worker{
		WorkerID:   workerID,
		WorkerPool: workerPool,
		Outputs:    o,
		Input:      make(chan Payload),
		quit:       make(chan bool),
	}, nil
}

// Start start a worker for mirror，处理mirror请求的内容，并将处理好的格式发给另外的进程处理
func (w *Worker) Start() {
	log.Printf("Mirror worker %s start!\n", w.WorkerID)
	go func() {
		for {
			// log.Printf("Worker %s start wait for payload.\n", w.WorkerID)
			w.WorkerPool <- w.Input
			// log.Printf("Worker %s wait for payload.\n", w.WorkerID)
			select {
			case payload := <-w.Input:
				var r Request
				r.Time = payload.Time
				//
				for header, value := range payload.Headers {
					switch header {
					case "X-Real-Ip":
						r.XRealIP = value[0]
					case "X-Original-Uri":
						// 去掉参数
						r.XOriginalURI = strings.SplitN(value[0], "?", 2)[0]
					case "X-Forwarded-For":
						r.XForwardFor = value
					case "X-Forwarded-Proto":
						r.XForwarderProto = value[0]
					case "Nginx-Ip":
						r.NginxIP = value[0]
					case "User-Agent":
						r.UserAgent = value
					}
				}
				r.Method = payload.Method

				// 另开一个线程把结果发给各个输出
				go func() {
					for _, output := range w.Outputs {
						// go counter.Add(r.XOriginalURI)
						output <- r
					}
				}()

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
