package count

import (
	"nginx_mirror/mirror"
	"testing"
	"time"
)

func TestDistpatchRun(t *testing.T) {
	begin := time.Now().Second()
	tickStart := time.NewTicker(time.Millisecond)

	var result CounterResult
	tpsChannel := make(chan mirror.Request)
	totalChannel := make(chan mirror.Request)
	dispatcher, err := NewDispatcher()
	if err != nil {
		t.Errorf("Error happend when new dispatcher, err: %s\n", err.Error())
		t.FailNow()
	}
	req1 := mirror.Request{
		XOriginalURI: "/web1",
	}
	req2 := mirror.Request{
		XOriginalURI: "/web2",
	}
	dispatcher.Run([]chan mirror.Request{tpsChannel, totalChannel})

	// 等待到下一秒开始
	for {
		<-tickStart.C
		if time.Now().Second() != begin {

			break
		}
	}
	tickStart.Stop()

	tick := time.NewTicker(time.Second)
	for i := 0; i < 100; i++ {
		tpsChannel <- req1
	}
	for i := 0; i < 50; i++ {
		tpsChannel <- req2
		totalChannel <- req2
	}
	time.Sleep(time.Millisecond)
	result = dispatcher.TPSGroup.Count.Read("")
	if result["/web1"] != 0 && result["/web2"] != 0 {
		t.Errorf("tps result: %+v\n", result)
		t.Fail()
	}
	result = dispatcher.TotalGroup.Count.Read("")
	if result["/web2"] != 50 {
		t.Errorf("tps result: %+v\n", result)
		t.Fail()
	}
	<-tick.C
	time.Sleep(500 * time.Millisecond)

	result = dispatcher.TPSGroup.Count.Read("")
	if result["/web1"] != 100 && result["/web2"] != 50 {
		t.Errorf("tps result: %+v\n", result)
		t.Fail()
	}
	<-tick.C
	time.Sleep(500 * time.Millisecond)
	result = dispatcher.TPSGroup.Count.Read("")
	if result["/web1"] != 0 && result["/web2"] != 0 {
		t.Errorf("tps result: %+v\n", result)
		t.Fail()
	}
	tick.Stop()

	dispatcher.Stop()
}
