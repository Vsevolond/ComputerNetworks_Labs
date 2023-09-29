package main

import (
	"flag"
	filedriver "github.com/goftp/file-driver"
	"github.com/goftp/server"
	"log"
)

func main() {
	root := "/Users/vsevolod/Documents/GO_projects/ComputerNetworks/UniServer/files"
	user := "user"
	pass := "1234"
	port := 2222
	host := "127.0.0.1"
	flag.Parse()
	factory := &filedriver.FileDriverFactory{
		RootPath: root,
		Perm:     server.NewSimplePerm("user", "group"),
	}
	opts := &server.ServerOpts{
		Factory:  factory,
		Port:     port,
		Hostname: host,
		Auth: &server.SimpleAuth{Name: user,
			Password: pass},
	}
	server := server.NewServer(opts)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}

//ftp user@localhost -p 2121
