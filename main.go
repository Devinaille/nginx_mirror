package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"nginx_mirror/config"
	"nginx_mirror/count"
	"nginx_mirror/mirror"
	"os"
	"time"
)

var mirrorConfig config.MirrorConfig

// 定义各种隧道
var payloadQueue chan mirror.Payload
var tpsRequest chan mirror.Request
var totalRequest chan mirror.Request

var mirrorDispatcher *mirror.Dispatcher
var counterDispatcher *count.Dispatcher

func init() {
	mirrorConfig = config.NewMirrorConfig()
	if len(os.Args) >= 2 {
		mirrorConfig.Load(os.Args[1])
	}
	log.Printf("URI: %s, port: %d, host: %s\n", mirrorConfig.MirrorURI, mirrorConfig.Port, mirrorConfig.Host)

	payloadQueue = make(chan mirror.Payload, 10000)
	tpsRequest = make(chan mirror.Request, 10000)
	totalRequest = make(chan mirror.Request, 10000)

}

// mirrorHandler 接收镜像过来的请求
func mirrorHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("get one!")
	payloadQueue <- mirror.Payload{
		Time:    time.Now(),
		Headers: r.Header,
		Method:  r.Method,
	}

	w.WriteHeader(http.StatusOK)
}

func tpsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		log.Printf("解析参数错误：%s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var cnt count.CounterResult
	if values, ok := r.Form["url"]; ok && len(values) > 0 {
		cnt = counterDispatcher.TPSGroup.Count.Read(string(values[0]))
	} else {
		cnt = counterDispatcher.TPSGroup.Count.Read("")
	}

	result, err := json.Marshal(cnt)
	if err != nil {
		log.Printf("查询url计数出错，%s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// w.WriteHeader(http.StatusOK)
	w.Write(result)
	return
}

func totalHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		log.Printf("解析参数错误：%s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var cnt count.CounterResult
	if values, ok := r.Form["url"]; ok && len(values) > 0 {
		cnt = counterDispatcher.TotalGroup.Count.Read(string(values[0]))
	} else {
		cnt = counterDispatcher.TotalGroup.Count.Read("")
	}

	result, err := json.Marshal(cnt)
	if err != nil {
		log.Printf("查询url计数出错，%s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// w.WriteHeader(http.StatusOK)
	w.Write(result)
	return
}

func initHTTPServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/mirror", mirrorHandler)
	mux.HandleFunc("/tps", tpsHandler)
	mux.HandleFunc("/total", totalHandler)
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", mirrorConfig.Host, mirrorConfig.Port),
		Handler: mux,
	}
	return server
}

func main() {
	var err error
	requestChannls := []chan mirror.Request{
		tpsRequest,
		totalRequest,
	}
	mirrorDispatcher, err = mirror.NewDispatcher(5)
	if err != nil {
		log.Fatalln(err)
	}
	counterDispatcher, err = count.NewDispatcher()
	if err != nil {
		log.Fatalln(err)
	}

	counterDispatcher.Run(requestChannls)
	mirrorDispatcher.Run(payloadQueue, requestChannls)
	server := initHTTPServer()

	server.ListenAndServe()
}
