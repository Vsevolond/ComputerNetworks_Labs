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
	"log"
	"math"
	"net/http"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		var req proto.Request
		err := c.ReadJSON(&req)
		if err != nil {
			log.Println("read:", err)
			return
		}
		switch req.Command {
		case "ok":
			fmt.Printf("ok\n")
		case "failed":
			if req.Data == nil {
				fmt.Printf("error: data field is absent in response\n")
			} else {
				var errorMsg string
				if err := json.Unmarshal(*req.Data, &errorMsg); err != nil {
					fmt.Printf("error: malformed data field in response\n")
				} else {
					fmt.Printf("failed: %s\n", errorMsg)
				}
			}
		case "get":
			if req.Data == nil {
				fmt.Printf("error: data field is absent in response\n")
			} else {
				var data proto.Circles
				if err := json.Unmarshal(*req.Data, &data); err != nil {
					fmt.Println("error: malformed data field in response, ", err)
				} else {
					intersections := intersect(data.Circle1, data.Circle2)
					send_response(c, "result", intersections)
				}
			}
		default:
			fmt.Printf("error: server reports unknown status %q\n", req.Command)
		}
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func send_response(conn *websocket.Conn, status string, data interface{}) {
	var raw json.RawMessage
	raw, _ = json.Marshal(&data)
	err := conn.WriteJSON(&proto.Response{status, &raw})
	if err != nil {
		log.Println("write:", err)
		return
	}
}

func intersect(circle1 proto.Circle, circle2 proto.Circle) []proto.Coord {
	x1, y1, x2, y2, r1, r2 := circle1.Center.X, circle1.Center.Y, circle2.Center.X, circle2.Center.Y,
		circle1.Radius, circle2.Radius
	if r1 > r2 {
		x1, y1, r1, x2, y2, r2 = x2, y2, r2, x1, y1, r1
	}
	if x1 == x2 && y1 == y2 && r1 == r2 {
		return []proto.Coord{{math.Inf(1), math.Inf(1)}}
	} else {

		d := math.Sqrt((x1-x2)*(x1-x2) + (y1-y2)*(y1-y2))
		if d > r1+r2 {
			return []proto.Coord{} //окружности не пересекаются
		}

		a := (r1*r1 - r2*r2 + d*d) / (2 * d)
		h := math.Sqrt(r1*r1 - a*a)

		x0 := x1 + a*(x2-x1)/d
		y0 := y1 + a*(y2-y1)/d
		var first proto.Coord
		first.X = x0 + h*(y2-y1)/d
		first.Y = y0 - h*(x2-x1)/d

		if a == r1 {
			return []proto.Coord{first}
		} //окружности соприкасаются
		var second proto.Coord
		second.X = x0 - h*(y2-y1)/d
		second.Y = y0 + h*(x2-x1)/d

		return []proto.Coord{first, second} //окружности пересекаются
	}
}
