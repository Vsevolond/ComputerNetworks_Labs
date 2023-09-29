package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var addr = flag.String("addr", "localhost:8030", "http service address")

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func check() string {
	resp, err := http.Get("https://bmstu.ru/")

	if err != nil {
		return "BAUMAN is not available"
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "BAUMAN is not available"
	}

	return "BAUMAN is norm"
}

func task3(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)

	for {
		msg := check()
		conn.WriteMessage(websocket.TextMessage, []byte(msg))
		time.Sleep(time.Second * 5)
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/task3", task3)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
