package main

import (
	"flag"
	"runtime"

	"websocket/internal/rdisplay"

	"log"
	"net/http"
	"time"
)

var (
	runStatus         bool
	path, port, frame string
)

func close(w http.ResponseWriter, r *http.Request) {
	if Encoder != nil && runStatus {
		runStatus = false
		log.Println("http close", time.Now().UnixNano(), runStatus)
	}
	w.Write([]byte("close"))
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func main() {
	flag.StringVar(&port, "port", "8020", "http port 8020")
	flag.StringVar(&frame, "frame", "10", "frame default 10")
	flag.StringVar(&path, "path", "/data/local/tmp/h264mini", "h264 file path")
	flag.Parse()

	defer func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("[Recovery] panic recovered: %v\n%s%s]", r, buf, string([]byte{27, 91, 48, 109}))
		}
	}()

	log.Println("init start", port, frame, path)
	http.HandleFunc("/api/start", wsNormalH264)
	http.Handle("/api/stop", http.HandlerFunc(close))
	http.Handle("/ping", http.HandlerFunc(ping))
	http.Handle("/public", http.FileServer(http.Dir("./public")))

	go rdisplay.InitCrontab(String2Int(frame))
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
