package count

import (
	"sync"
)

type cnt struct {
	rw  *sync.RWMutex
	cnt int64
}

// TotalCounter url计数器
type TotalCounter struct {
	rw    *sync.RWMutex
	count map[string]*cnt
}

// NewTotalCounter 初始化
func NewTotalCounter() *TotalCounter {
	return &TotalCounter{
		rw:    new(sync.RWMutex),
		count: make(map[string]*cnt),
	}
}

// Add 计数加1
// url string 需要计数的url
// t time.Time 未使用，兼容counter 用
func (t *TotalCounter) Add(url string) {
	// TODO: 只记录当天的
	var c *cnt
	var ok bool
	t.rw.RLock()

	// 判断URL有没有记录过，如果没有记录过则初始化，如果有则加1
	if c, ok = t.count[url]; !ok {
		t.rw.RUnlock()

		t.rw.Lock()
		if c, ok = t.count[url]; !ok {
			// 防止别的线程已经创建了
			c = &cnt{
				rw:  new(sync.RWMutex),
				cnt: 1,
			}
		}
		t.count[url] = c
		t.rw.Unlock()

	} else {
		t.rw.RUnlock()
		// .rw.Lock()
		c.rw.Lock()
		c.cnt = c.cnt + 1
		c.rw.Unlock()

	}
}

// Read 读取某一URL的值
func (t *TotalCounter) Read(url string) (cr CounterResult) {
	var c *cnt
	var ok bool
	cr = make(CounterResult)

	t.rw.RLock()

	// 返回具体某一个url的计数
	if url != "" {
		if c, ok = t.count[url]; ok {
			t.rw.RUnlock()
			c.rw.RLock()
			cr[url] = c.cnt
			c.rw.RUnlock()
			return
		}

		// 不存在
		t.rw.RUnlock()
		cr[url] = 0
		return
	}

	// 如果url参数为"",返回全部url的计数
	for url, c := range t.count {
		c.rw.RLock()
		cr[url] = c.cnt
		c.rw.RUnlock()
	}
	t.rw.RUnlock()
	return
}

// Reset 重置今天的数据
func (t *TotalCounter) Reset() {
	t.rw.Lock()
	t.count = make(map[string]*cnt)
	t.rw.Unlock()
}
