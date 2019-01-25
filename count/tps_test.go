package count

import (
	"testing"
)

func TestTpsAdd(t *testing.T) {
	url := ""
	tps := NewTPSCounter()
	tps.Add(url)
	result := tps.Read(url)
	if result[url] != 0 {
		t.Logf("result %+v\n", result[url])
		t.Fail()
	}
	tps.Reset()
	result = tps.Read(url)

	if result[url] != 1 {
		t.Logf("result %+v\n", result[url])
		t.Fail()
	}
}

func TestTpsReset(t *testing.T) {
	url := "/"
	tps := NewTPSCounter()
	tps.Add(url)
	result := tps.Read(url)

	if result[url] != 0 {
		t.Logf("result %+v\n", result[url])
		t.Fail()
	}

	tps.Reset()
	result = tps.Read(url)
	if result[url] != 1 {
		t.Logf("result %+v\n", result[url])
		t.Fail()
	}

	result = tps.Read(url)
	if result[url] != 1 {
		t.Logf("result %+v\n", result[url])
		t.Fail()
	}

	tps.Reset()
	result = tps.Read(url)
	if result[url] != 0 {
		t.Logf("result %+v\n", result[url])
		t.Fail()
	}
}

func TestTPSRead(t *testing.T) {
	url1 := "/foo"
	url2 := "/bar"
	tps := NewTPSCounter()
	tps.Add(url1)
	tps.Add(url2)
	tps.Add(url2)
	result1 := tps.Read(url1)
	if result1[url1] != 0 {
		t.Logf("result %d\n", result1[url1])
		t.Fail()
	}
	result2 := tps.Read(url2)
	if result2[url2] != 0 {
		t.Logf("result %d\n", result2[url2])
		t.Fail()
	}

	result3 := tps.Read("")
	if result3[url1] != 0 && result3[url2] != 0 {
		t.Logf("url1: %d, url2: %d\n", result3[url1], result3[url2])
		t.Fail()
	}

	tps.Reset()
	result1 = tps.Read(url1)
	if result1[url1] != 1 {
		t.Logf("result %d\n", result1[url1])
		t.Fail()
	}
	result2 = tps.Read(url2)
	if result2[url2] != 2 {
		t.Logf("result %d\n", result2[url2])
		t.Fail()
	}

	result3 = tps.Read("")
	if result3[url1] != 1 && result3[url2] != 1 {
		t.Logf("url1: %d, url2: %d\n", result3[url1], result3[url2])
		t.Fail()
	}
}
