package count

import (
	"nginx_mirror/mirror"
	"testing"
	"time"
)

func TestTPSWorker(t *testing.T) {
	var result CounterResult
	workerPool := make(chan chan mirror.Request, 1)
	counter := NewTPSCounter()
	request := mirror.Request{
		XOriginalURI: "/web",
	}
	w, err := NewWorker(workerPool, counter)
	if err != nil {
		t.Errorf("New tps worker failed, %s", err.Error())
		t.FailNow()
	}

	w.Start()
	if !w.Status() {
		t.Errorf("worker 启动失败。\n")
		t.FailNow()
	}
	queue := <-workerPool
	queue <- request

	// 等待消耗队列结束
	time.Sleep(time.Millisecond)

	result = counter.Read(request.XOriginalURI)
	if result[request.XOriginalURI] != 0 {
		t.Errorf("%s的计数不为0，为%d\n", request.XOriginalURI, result[request.XOriginalURI])
		t.Fail()
	}
	counter.Reset()
	result = counter.Read(request.XOriginalURI)
	if result[request.XOriginalURI] != 1 {
		t.Errorf("%s的计数不为1，为%d\n", request.XOriginalURI, result[request.XOriginalURI])
		t.Fail()
	}

	for i := 0; i < 10; i++ {
		queue := <-workerPool
		queue <- request
	}
	// 等待消耗队列结束
	time.Sleep(time.Millisecond)
	counter.Reset()
	result = counter.Read(request.XOriginalURI)
	if result[request.XOriginalURI] != 10 {
		t.Errorf("%s的计数不为10，为%d\n", request.XOriginalURI, result[request.XOriginalURI])
		t.Fail()
	}

	w.Stop()

	// 等待worker结束
	time.Sleep(time.Millisecond)
	if w.Status() {
		t.Errorf("worker 停止失败。状态为%t\n", w.Status())
		t.FailNow()
	}
}

func TestGetWorkerID(t *testing.T) {
	workerPool := make(chan chan mirror.Request, 1)
	counter := NewTPSCounter()
	w, err := NewWorker(workerPool, counter)
	if err != nil {
		t.Errorf("New tps worker failed, %s", err.Error())
		t.FailNow()
	}
	if w.ID() == "" {
		t.Errorf("Get worker ID failed, ID: %s\n", w.ID())
	}
}
