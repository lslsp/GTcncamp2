package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"strings"
	//"github.com/golang/glog"
)

func main() {
	//flag.Set("v", "4")
	//glog.V(2).Info("Starting http server...")
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/healthz", healthzHandler)
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	err := http.ListenAndServe(":80", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("entering root handler")
	user := r.URL.Query().Get("user")
	if user != "" {
		io.WriteString(w, fmt.Sprintf("hello [%s]\n", user))
	} else {
		io.WriteString(w, "hello [stranger]\n")
	}
	// 1.接收客户端request，并将requst中带的header写入reponse header
	io.WriteString(w, "===================Details of the http request header:============\n")
	for k, v := range r.Header {
		for _, vv := range v {
			w.Header().Set(k, vv)
			io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
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

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	//4.当访问localhost/healthz时，应返回200
	log.Printf("entering healthz handler")
	io.WriteString(w, "The HttpServer is Ready.\n")
}
