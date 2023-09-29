package main

import (
	"flag"
	"log"
	"net/http"
)

var addr1 = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/", db)
	log.Fatal(http.ListenAndServe(*addr1, nil))
}

func db(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}
