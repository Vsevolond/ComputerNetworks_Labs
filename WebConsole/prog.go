package main

import (
	"database/sql"
	"fmt"
	"log"
)

//func main() {
//	st := http.FileServer(http.Dir("static"))
//	http.Handle("/static/", http.StripPrefix("/static", st))
//	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//		url := "templates/login.html" //Страница входа
//		fmt.Println(r.Header["Cookie"][0])
//		if len(r.Header["Cookie"]) != 0 && r.Header["Cookie"][0] == "csrftoken=HEwDsMDzOVkobG6rwomRL1tTW5XMKq6F" {
//			url = "templates/base.html" //Страница после успешной авторизации
//		}
//		t, _ := template.ParseFiles(url)
//		t.Execute(w, "")
//	})
//	http.Handle("/js/", http.FileServer(http.Dir("templates")))
//	http.ListenAndServe(":8000", nil)
//}

func main() {
	database, err := sql.Open("sqlite3", "WebConsole/database.db")
	if err != nil {
		log.Fatal(err)
	}
	users_table := `CREATE TABLE users (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "Login" TEXT UNIQUE,
        "Password" TEXT);`
	query, err := database.Prepare(users_table)
	if err != nil {
		log.Fatal(err)
	}
	query.Exec()
	fmt.Println("Table created successfully!")
}
