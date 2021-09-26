package main

import (
	"fmt"
	"k8s.io/klog/v2"
	"net/http"
	"os"
	"strings"
)

func main() {
	// 设置路由
	http.HandleFunc("/healthz", health)
	// 设置监听的端口
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		klog.Fatalln("ListenAndServe: ", err)
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	err := writeHeader(w, r)

	// 判断是否是get请求
	var code int
	var msg string
	if strings.ToLower(r.Method) != "get" {
		code = 405
		msg = "Method not allowed"
	} else {
		code = 200
		msg = "OK"
	}

	writeBody(code, msg, w, err)

	klog.Infoln("Method:", r.Method, "Url:", r.URL, "StatusCode:", code)
}

// 写入消息体
func writeBody(statusCode int, msg string, w http.ResponseWriter, headerErr error) {
	if headerErr != nil {
		statusCode = 500
		msg = "Internal server error"
	}

	w.WriteHeader(statusCode)
	buffer := []byte(msg)
	write, writeErr := w.Write(buffer)
	if writeErr != nil {
		klog.Errorln("code:", write, "error:", writeErr)
	}
}

// RequestHeader写入ResponseHeader
func writeHeader(w http.ResponseWriter, r *http.Request) error {
	if w == nil {
		err := fmt.Errorf("http.ResponseWriter is nil")
		klog.Errorln("responseWriter is nil,", err)
		return err
	}

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

	return nil
}
