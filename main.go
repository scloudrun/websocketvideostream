package main

import (
	"flag"
	"io/ioutil"
	"path/filepath"
	"runtime"

	"golang.org/x/net/websocket"

	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var path, port, frame string

func wsH264(ws *websocket.Conn) {
	Init()
	num := 1
	diff := 1000 / String2Int(frame)
	defer ws.Close()
	for {
		time.Sleep(time.Duration(diff) * time.Millisecond)

		if num == 20 {
			initEncoder()
			num = 0
		}

		files := FileWalk(path)
		if len(files) >= 2 {
			file := files[len(files)-2]
			fileByte, err := getEncode(file)
			if err != nil || len(fileByte) == 0 {
				continue
			}
			err = websocket.Message.Send(ws, fileByte)
			if err != nil {
				log.Println(err)
				break
			}
		}
	}
	log.Println("send over socket\n")
	if Encoder != nil {
		Encoder.Close()
	}
}

func close(w http.ResponseWriter, r *http.Request) {
	if Encoder != nil {
		Encoder.Close()
	}
	w.Write([]byte("close"))
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
	http.Handle("/api/wsh264", websocket.Handler(wsH264))
	http.HandleFunc("/api/start", wsNormalH264)
	http.Handle("/api/stop", http.HandlerFunc(close))
	http.Handle("/", http.FileServer(http.Dir("./public")))

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

// GetFileByte def
func GetFileByte(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil && os.IsNotExist(err) {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

// FileWalk def
func FileWalk(fileDir string) []string {
	start, err := os.Stat(fileDir)
	if err != nil || !start.IsDir() {
		return nil
	}
	var targets []string
	filepath.Walk(fileDir, func(fpath string, fi os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !fi.Mode().IsRegular() {
			return nil
		}
		targets = append(targets, fpath)
		return nil
	})
	return targets
}

// String2Int def
func String2Int(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}
