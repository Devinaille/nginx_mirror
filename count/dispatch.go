package count

import (
	"errors"
	"fmt"
	"log"
	"nginx_mirror/mirror"

	"github.com/robfig/cron"
)

// dispatchGroup 存放调度组的基本信息
type dispatchGroup struct {
	Name       string
	MaxWorkers int
	WorkerPool chan chan mirror.Request
	Workers    []*Worker
	Count      Counter
	quit       chan bool
}

// Dispatcher 调度器，用于调用counter的worker
type Dispatcher struct {
	TPSGroup   *dispatchGroup
	TotalGroup *dispatchGroup
	c          *cron.Cron
	quit       chan bool
}

// newDispatchGroup 初始化调度组
func newDispatchGroup(name string, maxWorkers int, c Counter) (*dispatchGroup, error) {

	if maxWorkers == 0 {
		maxWorkers = 5
	}
	if name == "" {
		name = "Group"
	}
	return &dispatchGroup{
		Name:       name,
		MaxWorkers: maxWorkers,
		WorkerPool: make(chan chan mirror.Request, maxWorkers),
		Workers:    make([]*Worker, maxWorkers),
		Count:      c,
		quit:       make(chan bool),
	}, nil
}

// NewDispatcher 初始化counter调度器
func NewDispatcher() (*Dispatcher, error) {
	tpsGroup, err := newDispatchGroup("TPS GROUP", 5, NewTPSCounter()) //TODO: 以后从配置文件读取
	if err != nil {
		return nil, errors.New("Tps Group init failed")
	}
	totalGroup, err := newDispatchGroup("TOTAL GROUP", 5, NewTotalCounter()) //TODO: 以后从配置文件读取
	if err != nil {
		return nil, errors.New("Total Group init failed")
	}
	return &Dispatcher{
		TPSGroup:   tpsGroup,
		TotalGroup: totalGroup,
		c:          cron.New(),
		quit:       make(chan bool),
	}, nil
}

// Run 启动调度组
func (d *dispatchGroup) Run(ch chan mirror.Request) {
	for i := 0; i < d.MaxWorkers; i++ {
		worker, err := NewWorker(d.WorkerPool, d.Count)
		if err != nil {
			log.Printf("One worker in tps group create failed, %s\n", err.Error())
			continue
		}
		d.Workers[i] = worker
		worker.Start()
	}
	go d.dispatch(ch)
}

func (d *dispatchGroup) dispatch(ch chan mirror.Request) {
	for {
		select {
		case request := <-ch:
			go func(r mirror.Request) {
				requestChannel := <-d.WorkerPool
				requestChannel <- r
			}(request)
		case <-d.quit:
			close(ch) // TODO: 也许要mirror那边关闭通道
			for _, w := range d.Workers {
				fmt.Printf("%p\n", w.quit)
				w.Stop()
			}
			return
		}
	}
}

func (d *dispatchGroup) Stop() {
	log.Printf("Stop dispatcher group %s!", d.Name)
	d.quit <- true
}

// Run 启用counter相关调度器
// inputs 来自mirror workers的outputs
func (d *Dispatcher) Run(inputs []chan mirror.Request) {

	d.TPSGroup.Run(inputs[0])
	d.TotalGroup.Run(inputs[1])

	// 启动tps、total的定时重置任务
	d.c.AddFunc("* * * * * *", d.TPSGroup.Count.Reset)
	d.c.AddFunc("0 0 0 * * *", d.TotalGroup.Count.Reset)
	go d.c.Start()

}

// Stop 停用counter相关调度器
// inputs 来自mirror workers的outputs
func (d *Dispatcher) Stop() {
	d.c.Stop()
	d.TPSGroup.Stop()
	d.TotalGroup.Stop()
}
