package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"strings"
	"time"

	//模块十 作业2. 为 HTTPServer 项目添加延时 Metric (1/3)
	"github.com/lslsp/GTcncamp2/service2/metrics"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	//使用纳秒时间戳,每次启动程序的时候随机数不一样
	rand.Seed(time.Now().UTC().UnixNano())

	log.Printf("Starting service2...")

	//模块十 作业2. 为 HTTPServer 项目添加延时 Metric (2/3)
	metrics.Register()

	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/healthz", healthzHandler)
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	//模块十 作业2. 为 HTTPServer 项目添加延时 Metric (3/3)
	mux.Handle("/metrics", promhttp.Handler())

	if err := http.ListenAndServe(":80", mux); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}

}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("entering root handler")

	timer := metrics.NewTimer()
	defer timer.ObserveTotal()
	//模块十 作业1. 为HTTPServer添加0-2秒的随机延时
	delay := rand.Intn(2000)
	time.Sleep(time.Millisecond * time.Duration(delay))
	io.WriteString(w, fmt.Sprintf("<h1>%d</h1>", delay))

	//模块十二 作业3. 考虑Opentracing的接入（3/3）
	user := r.URL.Query().Get("user")
	if user != "" {
		io.WriteString(w, fmt.Sprintf("hello %s !<br/>", user))
	} else {
		io.WriteString(w, "hello [stranger]<br/>\n")
	}
	//模块二 作业1. 接收客户端request，并将requst中带的header写入reponse header
	io.WriteString(w, "===================Details of the http request header:============<br>")
	for k, v := range r.Header {
		for _, vv := range v {
			w.Header().Set(k, vv)
			io.WriteString(w, fmt.Sprintf("%s=%s<br>", k, v))
		}
	}

	//模块二 作业2. 读取当前系统的环境变量中的VERSION配置，并写入response header
	version := os.Getenv("VERSION")
	w.Header().Set("VERSION", version)

	//模块二 作业3. Server端记录访问日志（包括客户端IP，Http返回码），输出到server端的标准输出
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
	//模块二 作业4. 当访问localhost/healthz时，应返回200
	log.Printf("entering healthz handler")
	io.WriteString(w, "<h1>Ready</h1>")
}
