package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/healthz", healthz)
	err := http.ListenAndServe(":80", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	// 1.接收客户端request，并将requst中带的header写入reponse header
	for k, v := range r.Header {
		for _, vv := range v {
			w.Header().Set(k, vv)
		}
	}

	// 2.读取当前系统的环境变量中的VERSION配置，并写入response header
	version := os.Getenv("VERSION")
	w.Header().Set("VERSION", version)

	//3.Server端记录访问日志（包括客户端IP，Http返回码），输出到server端的标准输出
	clientip := ClientIP(r)
	log.Printf("clientip: %s", clientip)
	log.Printf("Response code: %d", 200)
}

func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

func healthz(w http.ResponseWriter, r *http.Request) {
	//4.当访问localhost/healthz时，应返回200
	fmt.Fprintf(w, "200")
}
