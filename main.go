package main

import (
	"fmt"
	"log"
	"net/http"
	mc "nginx_mirror/config"
	"nginx_mirror/mirror"
	"os"
)

// mirrorConfig 配置文件全局变量
var mirrorConfig mc.MirrorConfig
var payloadQueue chan mirror.Payload
var dispatcher *mirror.Dispatcher

// mirrorHandler 接收镜像过来的请求
func mirrorHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("get one!")
	payloadQueue <- mirror.Payload{
		Headers: r.Header,
		Method:  r.Method,
	}

	w.WriteHeader(http.StatusOK)
}

// mirrorHandler 获取TPS
func getTPSHandler(w http.ResponseWriter, r *http.Request) {

}

// 其他接口

func init() {
	mirrorConfig = mc.NewMirrorConfig()
	if len(os.Args) >= 2 {
		mirrorConfig.Load(os.Args[1])
	}
	log.Printf("URI: %s, port: %d, host: %s\n", mirrorConfig.MirrorURI, mirrorConfig.Port, mirrorConfig.Host)

}

func main() {
	var err error
	payloadQueue = make(chan mirror.Payload, 10000) // TODO: 队列长度可修改
	dispatcher, err = mirror.NewDispatcher(1)       //TODO: worker个数可修改，可以使用环境变量
	if err != nil {
		log.Fatalln(err)
	}

	dispatcher.Run(payloadQueue)

	mux := http.NewServeMux()

	// stop all process
	// go func() {
	// 	// close(payloadQueue)
	// 	time.Sleep(time.Second * 5)
	// 	dispatcher.Stop()
	// }()
	mux.HandleFunc(mirrorConfig.MirrorURI, mirrorHandler)
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", mirrorConfig.Host, mirrorConfig.Port), mux)
	if err != nil {
		log.Fatalln(err)
	}
	// TODO: 优雅的退出httpserver
}
