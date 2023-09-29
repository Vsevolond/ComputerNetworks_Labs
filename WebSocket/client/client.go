// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"ComputerNetworks/WebSocket/proto"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/skorobogatov/input"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
)

func main() {
	var addr = flag.String("addr", "localhost:8080", "http service address")
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	go func() {
		for {
			var resp proto.Response
			err := c.ReadJSON(&resp)
			if err != nil {
				log.Println("read:", err)
				return
			}
			switch resp.Status {
			case "ok":
				fmt.Printf("ok\n")
			case "failed":
				if resp.Data == nil {
					fmt.Printf("error: data field is absent in response\n")
				} else {
					var errorMsg string
					if err := json.Unmarshal(*resp.Data, &errorMsg); err != nil {
						fmt.Printf("error: malformed data field in response\n")
					} else {
						fmt.Printf("failed: %s\n", errorMsg)
					}
				}
			case "result":
				if resp.Data == nil {
					//fmt.Printf("error: data field is absent in response\n")
					log.Println("circles are equal")
				} else {
					var res []proto.Coord
					if err := json.Unmarshal(*resp.Data, &res); err != nil {
						fmt.Printf("error: malformed data field in response\n")
					} else {
						fmt.Println(res)
					}
				}
			default:
				fmt.Printf("error: server reports unknown status %q\n", resp.Status)
			}
		}
	}()

	for {
		str := input.Gets()
		switch str {
		case "exit":
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("close:", err)
				return
			}
			return
		case "get":
			log.Println("# First Circle:")
			fmt.Print("Center = ")
			center1X, center1Y := center(input.Gets())
			fmt.Print("Radius = ")
			radius1, _ := strconv.Atoi(input.Gets())
			log.Println("# Second Circle:")
			fmt.Print("Center = ")
			center2X, center2Y := center(input.Gets())
			fmt.Print("Radius = ")
			radius2, _ := strconv.Atoi(input.Gets())
			req := &proto.Circles{
				Circle1: proto.Circle{
					Center: proto.Coord{X: center1X, Y: center1Y},
					Radius: float64(radius1),
				},
				Circle2: proto.Circle{
					Center: proto.Coord{X: center2X, Y: center2Y},
					Radius: float64(radius2),
				},
			}
			send_request(c, "get", &req)
		default:
			log.Println("ERR: unknown command")
		}
	}
}

func center(str string) (float64, float64) {
	coord := strings.Split(str, " ")
	a, b := coord[0], coord[1]
	x, err := strconv.Atoi(a)
	if err != nil {
		log.Println("can't convert")
	}
	y, err := strconv.Atoi(b)
	if err != nil {
		log.Println("can't convert")
	}
	return float64(x), float64(y)
}

func send_request(conn *websocket.Conn, command string, data interface{}) {
	var raw json.RawMessage
	raw, _ = json.Marshal(&data)
	err := conn.WriteJSON(&proto.Request{command, &raw})
	if err != nil {
		log.Println("write:", err)
		return
	}
}
