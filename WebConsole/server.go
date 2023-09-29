// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type User struct {
	login    string
	password string
}

var path, _ = os.Getwd()
var dir = strings.Split(path, "/")

var upgrader = websocket.Upgrader{} // use default options

func getPath(dir []string) string {
	path := ""
	for _, d := range dir {
		path += "/" + d
	}
	return path
}

func main() {

	st := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static", st))

	http.HandleFunc("/terminal", func(w http.ResponseWriter, r *http.Request) {
		// Upgrade upgrades the HTTP server connection to the WebSocket protocol.
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade failed: ", err)
			return
		}
		defer conn.Close()

		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read failed:", err)
				break
			}
			//fmt.Println(os.Getwd())
			//оs.Chdir()
			input := string(message)
			args := strings.Split(input, " ")
			cmd := exec.Command(args[0], args[1:]...)
			out, err := cmd.Output()
			if err != nil {
				out = []byte(err.Error())
			}
			if args[0] == "cd" {
				switch args[0] {
				case "..":
					dir = append(dir[:len(dir)-1])
				default:
					dir = append(dir, args[1])
				}
			}
			path = getPath(dir)
			err = os.Chdir(path)
			if err != nil {
				out = []byte(err.Error())
				dir = append(dir[:len(dir)-1])
				path = getPath(dir)
				os.Chdir(path)
			}
			err = conn.WriteMessage(mt, out)
			if err != nil {
				log.Println("write failed:", err)
				break
			}
		}
	})

	//http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
	//	//http.ServeFile(w, r, "templates/login.html")
	//	conn, err := upgrader.Upgrade(w, r, nil)
	//	if err != nil {
	//		log.Print("upgrade failed: ", err)
	//		return
	//	}
	//	defer conn.Close()
	//
	//	for {
	//		mt, message, err := conn.ReadMessage()
	//		if err != nil {
	//			fmt.Println("read failed:", err)
	//			break
	//		}
	//		user := make(map[string]string)
	//		err = json.Unmarshal(message, &user)
	//		if err != nil {
	//			fmt.Println(err)
	//			break
	//		} else {
	//			database, err := sql.Open("sqlite3", "database.db")
	//			if err == nil {
	//				password_db := database.QueryRow("SELECT password FROM users WHERE login=?", user["login"])
	//				var pass string
	//				if err = password_db.Scan(&pass); err != nil {
	//					if err == sql.ErrNoRows {
	//						conn.WriteMessage(mt, []byte("bad login"))
	//					} else {
	//						fmt.Println(err)
	//					}
	//				} else {
	//					if pass != user["password"] {
	//						conn.WriteMessage(mt, []byte("bad password"))
	//					} else {
	//						conn.WriteMessage(mt, []byte("ok"))
	//					}
	//				}
	//				database.Close()
	//				break
	//			}
	//		}
	//	}
	//})
	//
	//http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
	//	// Upgrade upgrades the HTTP server connection to the WebSocket protocol.
	//	conn, err := upgrader.Upgrade(w, r, nil)
	//	if err != nil {
	//		log.Print("upgrade failed: ", err)
	//		return
	//	}
	//	defer conn.Close()
	//
	//	for {
	//		mt, message, err := conn.ReadMessage()
	//		if err != nil {
	//			fmt.Println("read failed:", err)
	//			break
	//		}
	//		user := make(map[string]string)
	//		err = json.Unmarshal(message, &user)
	//		if err != nil {
	//			fmt.Println(err)
	//			break
	//		} else {
	//			database, err := sql.Open("sqlite3", "database.db")
	//			if err == nil {
	//				_, err = database.Exec("insert into users (login, password) values (?, ?)",
	//					user["login"], user["password"])
	//				if err != nil {
	//					conn.WriteMessage(mt, []byte("exist"))
	//				} else {
	//					conn.WriteMessage(mt, []byte("ok"))
	//				}
	//				database.Close()
	//				break
	//			}
	//		}
	//	}
	//})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/base.html")
	})

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		return
	}
}
