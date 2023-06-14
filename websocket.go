package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

func wsNormalH264(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	Init()
	num := 1
	diff := 1000 / String2Int(frame)
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
			err = c.WriteMessage(websocket.BinaryMessage, fileByte)
			if err != nil {
				log.Println("write:", err)
				break
			}
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
