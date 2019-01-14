package main

import (
	"fmt"
	"log"
	"net/http"
	mc "nginx_mirror/config"
	"os"
)

// mirrorConfig 配置文件全局变量
var mirrorConfig mc.MirrorConfig

// mirrorHandler 接收镜像过来的请求
func mirrorHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("get one")
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
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc(mirrorConfig.MirrorURI, mirrorHandler)
	err := http.ListenAndServe(fmt.Sprintf(":%d", mirrorConfig.Port), mux)
	if err != nil {
		log.Fatalln(err)
	}
}
