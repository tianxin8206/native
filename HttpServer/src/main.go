package main

import (
	"HttpServer/src/config"
	"HttpServer/src/metrics"
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func main() {
	//读取json
	configuration, err := config.LoadConfiguration()
	if err != nil {
		log.Fatalln("json读取错误")
		return
	}

	metrics.Register()

	// 设置路由
	http.HandleFunc("/healthz", health)
	http.HandleFunc("/hello", rootHandler)
	http.Handle("/metrics", promhttp.Handler())

	// 设置监听的端口
	server := &http.Server{Addr: ":" + strconv.Itoa(configuration.Port), Handler: nil}
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalln("ListenAndServe: ", err)
		}
	}()

	//优雅停止
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case <-sigs:
		log.Println("收到停止信号")
		cancelErr := server.Shutdown(context.Background())
		if cancelErr != nil {
			log.Fatalln("Shutdown: ", cancelErr)
		}
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	if w == nil {
		err := fmt.Errorf("http.ResponseWriter is nil")
		log.Fatalln("responseWriter is nil,", err)
		return
	}

	// 判断是否是get请求
	var (
		code int
		msg  string
	)
	if strings.ToLower(r.Method) != "get" {
		code = 405
		msg = "Method not allowed"
	} else {
		code = 200
		msg = "OK"
	}

	writeBody(code, msg, w, r)

	log.Println("Method:", r.Method, "Url:", r.URL, "StatusCode:", code)
}

// 写入消息体
func writeBody(statusCode int, msg string, w http.ResponseWriter, r *http.Request) {
	writeHeader(w, r)

	w.WriteHeader(statusCode)
	buffer := []byte(msg)
	write, writeErr := w.Write(buffer)
	if writeErr != nil {
		log.Println("code:", write, "error:", writeErr)
	}
}

// RequestHeader写入ResponseHeader
func writeHeader(w http.ResponseWriter, r *http.Request) {
	for header := range r.Header {
		// 跳过Content-Length Header
		if header == "Content-Length" {
			continue
		}

		values := r.Header[header]
		for index := range values {
			w.Header().Set(header, values[index])
		}
	}

	// 读取env
	version := os.Getenv("Version")
	w.Header().Set("Version", version)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("entering root handler")
	timer := metrics.NewTimer()
	defer timer.ObserveTotal()
	user := r.URL.Query().Get("user")
	delay := randInt(10, 2000)
	time.Sleep(time.Millisecond * time.Duration(delay))
	if user != "" {
		_, _ = io.WriteString(w, fmt.Sprintf("hello [%s]\n", user))
	} else {
		_, _ = io.WriteString(w, "hello [stranger]\n")
	}
	_, _ = io.WriteString(w, "===================Details of the http request header:============\n")
	for k, v := range r.Header {
		_, _ = io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
	}
	log.Println("Respond in %d ms", delay)
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}
