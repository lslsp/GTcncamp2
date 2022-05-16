package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/pprof"
	"strings"
	"time"

	//模块十 作业2. 为 HTTPServer 项目添加延时 Metric (1/3)
	"github.com/lslsp/GTcncamp2/service0/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	//使用纳秒时间戳,每次启动程序的时候随机数不一样
	rand.Seed(time.Now().UTC().UnixNano())

	log.Printf("Starting service0")

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

	io.WriteString(w, "===================Details of the http request header:============\n")

	//模块十二 作业3. 考虑Opentracing的接入（1/3）
	req, err := http.NewRequest("GET", "http://service1", nil)

	if err != nil {
		fmt.Printf("%s", err)
	}
	lowerCaseHeader := make(http.Header)
	for key, value := range r.Header {
		lowerCaseHeader[strings.ToLower(key)] = value
	}
	log.Printf("headers:", lowerCaseHeader)
	req.Header = lowerCaseHeader
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("HTTP get failed with error: ", "error", err)
	} else {
		log.Printf("HTTP get succeeded")
	}
	if resp != nil {
		resp.Write(w)
	}
	log.Printf("Respond in %d ms", delay)
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	//模块二 作业4. 当访问localhost/healthz时，应返回200
	log.Printf("entering healthz handler")
	io.WriteString(w, "<h1>Ready</h1>")

	for k, v := range r.Header {
		io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
	}
}
