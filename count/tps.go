package count

import (
	"sync"
)

// tps
// 	rw 读写锁
// 	centisecondCnt 按照毫秒计数
type tps struct {
	rw    *sync.RWMutex
	Count int64
}

// newTPS 初始化TPSCnt
func newTPS() *tps {
	return &tps{
		rw:    new(sync.RWMutex),
		Count: 0,
	}
}

func (t *tps) add() {
	t.Count++
}

// count 读取tps计数
func (t *tps) count() int64 {
	var cnt int64

	// cnt = 0
	t.rw.RLock()
	cnt = t.Count
	t.rw.RUnlock()
	return cnt
}

// TPSCounter 存储按URL分类的TPS
type TPSCounter struct {
	rw            *sync.RWMutex
	lastSecondTPS map[string]*tps
	tps           map[string]*tps
}

// NewTPSCounter 初始化TPSCount计数
func NewTPSCounter() *TPSCounter {
	return &TPSCounter{
		rw:            new(sync.RWMutex),
		lastSecondTPS: make(map[string]*tps),
		tps:           make(map[string]*tps),
	}
}

// Add 根据url计数加1，如果不存在则初始化，并将当前毫秒置一
func (t *TPSCounter) Add(url string) {
	var tc *tps
	var ok bool
	t.rw.RLock()
	// msec := getMsecPos(requestTime)
	if tc, ok = t.tps[url]; !ok {
		t.rw.RUnlock()
		t.rw.Lock()
		if tc, ok = t.tps[url]; !ok {
			tc = newTPS()
		}
		tc.add()
		t.tps[url] = tc
		t.rw.Unlock()

	} else {
		t.rw.RUnlock()

		tc.rw.Lock()
		tc.add()
		tc.rw.Unlock()
	}
}

func (t *TPSCounter) Read(url string) (cr CounterResult) {
	var lst *tps
	var ok bool
	cr = make(CounterResult)
	t.rw.RLock()
	if url != "" {
		if lst, ok = t.lastSecondTPS[url]; !ok {
			// URL不存在
			cr[url] = 0
			t.rw.RUnlock()
			return
		}
		cr[url] = lst.count()
		t.rw.RUnlock()
		return

	}

	// 返回所有URL的TPS
	for k, v := range t.lastSecondTPS {
		cr[k] = v.count()
	}
	t.rw.RUnlock()
	return
}

// Reset 将当前的TPS放到上一秒的TPS中，并重置当前秒的数据
func (t *TPSCounter) Reset() {
	t.rw.Lock()
	t.lastSecondTPS = t.tps
	t.tps = make(map[string]*tps)
	t.rw.Unlock()
}
