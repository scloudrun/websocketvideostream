package main

import (
	"log"
	"time"

	"net/http"

	"github.com/gorilla/websocket"
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
	if runStatus {
		log.Print("runStatus true can't start")
		return
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	if err = initEncoder(); err != nil {
		log.Print("initEncoder:", err)
		return
	}

	delta := time.Duration(1000/String2Int(frame)) * time.Millisecond
	runStatus = true
	num := 1
	log.Println("wsH264 start", runStatus)

	for {
		startedAt := time.Now()
		if !runStatus {
			log.Println("runStatus false break")
			Encoder.Close()
			break
		}

		if num == 20 {
			log.Println("initEncoders", num)
			Encoder.Close()
			initEncoder()
			num = 0
		}
		num++

		files := FileWalk(path)
		if len(files) >= 2 {
			file := files[len(files)-2]
			if !runStatus {
				continue
			}

			fileByte, err := getEncode(file)
			if err != nil || len(fileByte) == 0 {
				continue
			}

			err = c.WriteMessage(websocket.BinaryMessage, fileByte)
			if err != nil {
				log.Println("socket close write:", err)
				break
			}
		}

		ellapsed := time.Now().Sub(startedAt)
		sleepDuration := delta - ellapsed
		if sleepDuration > 0 {
			time.Sleep(sleepDuration)
		}

	}
	log.Println("wsH264 stop ,websocket stop", runStatus)
	time.Sleep(time.Second * 1)
}
