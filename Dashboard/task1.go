package main

import (
	"flag"
	"github.com/gorilla/websocket"
	ssh "github.com/helloyi/go-sshclient"
	"log"
	"net/http"
	"strings"
	"time"
)

var addr = flag.String("addr", "localhost:8010", "http service address")

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func task1(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()

	client, err := ssh.DialWithPasswd("151.248.113.144:443", "test", "SDHBCXdsedfs222")
	if err != nil {
		log.Println(err)
	}
	defer client.Close()

	for {
		out, err := client.Cmd(`find -maxdepth 1 -name "achtung.txt" -exec cat {} +`).Output()
		if err != nil {
			log.Println(err)
		}

		if strings.Compare(string(out), "") == 0 {
			err = conn.WriteMessage(websocket.TextMessage, []byte("norm"))
			if err != nil {
				log.Println("write:", err)
				break
			}
		} else {
			err = conn.WriteMessage(websocket.TextMessage, out)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}

		time.Sleep(2 * time.Second)
	}
}
func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/task1", task1)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

//echo some-text  > filename.txt
