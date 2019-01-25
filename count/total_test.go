package count

import (
	"testing"
)

func TestTotalAdd(t *testing.T) {
	url := "/"
	tps := NewTotalCounter()
	tps.Add(url)
	result := tps.Read(url)
	if result[url] != 1 {
		t.Fail()
	}
	tps.Add(url)
	result = tps.Read(url)
	if result[url] != 2 {
		t.Fail()
	}
}

func TestTotalReset(t *testing.T) {
	url := "/"
	tps := NewTotalCounter()
	tps.Add(url)
	result := tps.Read(url)
	if result[url] != 1 {
		t.Fail()
	}
	tps.Reset()
	result = tps.Read(url)
	if result[url] != 0 {
		t.Fail()
	}
}

func TestTotalRead(t *testing.T) {
	url1 := "/foo"
	url2 := "/bar"
	tps := NewTotalCounter()
	tps.Add(url1)
	tps.Add(url2)
	result1 := tps.Read(url1)
	if result1[url1] != 1 {
		t.Fail()
	}
	result2 := tps.Read(url2)
	if result2[url2] != 1 {
		t.Fail()
	}

	result3 := tps.Read("")
	if result3[url1] != 1 && result3[url2] != 1 {
		t.Fail()
	}
}
