package count

// CounterResult 结果集
type CounterResult map[string]int64

// Counter worker的输出接口，需要实现Add、Read方法
type Counter interface {
	Add(url string)
	Read(url string) CounterResult
	Reset()
}
