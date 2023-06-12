package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"golang.org/x/net/websocket"

	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var Mode string

func wsH264(ws *websocket.Conn) {
	fmt.Printf("new socket\n")
	num := 1
	for {
		time.Sleep(100 * time.Millisecond)

		if num == 20 {
			initEncoder()
			fmt.Println("new encoder")
			num = 0
		}

		if Mode == "local" {
			files := FileWalk("/data/local/tmp/h264img")
			for i := 0; i < 200; i++ {
				time.Sleep(100 * time.Millisecond)
				fmt.Println("here i", i)
				if len(files) >= 2 {
					file := files[len(files)-2]
					file = files[i]
					fileByte, err := getEncode(file)
					if err != nil || len(fileByte) == 0 {
						return
					}
					err = websocket.Message.Send(ws, fileByte)
					fmt.Println("here", file, err)
					if err != nil {
						log.Println(err)
						break
					}
				}
			}
		}

		if Mode != "local" {
			files := FileWalk("/data/local/tmp/h264mini")
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
	}

	log.Println("send over socket\n")
}

func wsflv(ws *websocket.Conn) {
	fmt.Printf("new socket\n")

	fi, err := os.Open("./test.flv")
	if err != nil {
		log.Fatal(err)
	}
	defer fi.Close()

	msg := make([]byte, 1024*5)

	for {
		time.Sleep(4 * time.Millisecond)

		/*lenNum, err := fi.Read(lenBytes)
		if (err != nil && err != io.EOF) || lenNum != 4 {
			log.Println(err)
			time.Sleep(1 * time.Second)
			break
		}*/

		var lenreal int32
		lenreal = 1024 * 5

		log.Println(lenreal)

		n, err := fi.Read(msg[0:lenreal])
		if (err != nil && err != io.EOF) || n != int(lenreal) {
			log.Println(err)
			time.Sleep(1 * time.Second)
			break
		}

		err = websocket.Message.Send(ws, msg[:n])
		if err != nil {
			log.Println(err)
			break
		}
	}

	log.Println("send over socket\n")
}

func wsMpeg1(ws *websocket.Conn) {
	fmt.Printf("new socket\n")

	buf := make([]byte, 10)
	buf[0] = 'j'
	buf[1] = 's'
	buf[2] = 'm'
	buf[3] = 'p'
	buf[4] = 0x01
	buf[5] = 0x40
	buf[6] = 0x0
	buf[7] = 0xf0
	websocket.Message.Send(ws, buf[:8])

	fi, err := os.Open("./test.mpeg")
	if err != nil {
		log.Fatal(err)
	}
	defer fi.Close()

	msg := make([]byte, 1024*1)
	for {
		time.Sleep(40 * time.Millisecond)
		n, err := fi.Read(msg)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		if 0 == n {
			time.Sleep(1 * time.Second)
			break
		}
		err = websocket.Message.Send(ws, msg[:n])
		if err != nil {
			log.Println(err)
			break
		}
	}
	fmt.Printf("send over socket\n")
}

func main() {
	var mode string
	flag.StringVar(&mode, "mode", "", "mode type local")
	flag.Parse()

	if mode == "" {
		Mode = "local"
	}

	fmt.Println("here start mode ", Mode)

	http.Handle("/wsh264", websocket.Handler(wsH264))
	http.Handle("/wsmpeg", websocket.Handler(wsMpeg1))
	http.Handle("/wsflv", websocket.Handler(wsflv))

	http.Handle("/", http.FileServer(http.Dir("./public")))

	err := http.ListenAndServe(":8088", nil)
	fmt.Println("here start")
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
	//for _, target := range targets {
	//log.Logger.Debug("RecoverFromFile fileWalk file[%s]", target)
	//}
	return targets
}
