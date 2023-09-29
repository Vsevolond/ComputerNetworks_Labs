package main

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"github.com/mmcdole/gofeed"
	"github.com/skorobogatov/input"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func interact(client *ftp.ServerConn) {
	exit := false
	dir := []string{"~"}
	for {
		fmt.Print(dir[len(dir)-1] + " % ")
		str := input.Gets()
		command, path := split(str)
		switch command {
		case "exit":
			exit = true
			if err := client.Quit(); err != nil {
				log.Fatal(err)
			}
		case "get":
			file, err := client.Retr(path)
			if err != nil {
				log.Fatal(err)
			}
			savedfile, err := os.Create("/Users/vsevolod/Downloads/" + path)
			if err != nil {
				log.Fatal(err)
			}
			_, err = io.Copy(savedfile, file)
			if err != nil {
				log.Fatal(err)
			}
			savedfile.Close()
			defer savedfile.Close()
			file.Close()
			defer file.Close()
		case "put":
			file, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}
			ind := strings.LastIndex(path, "/")
			name := path[ind+1:]
			err = client.Stor(name, file)
			if err != nil {
				log.Fatal(err)
			}
			file.Close()
			defer file.Close()
		case "mkdir":
			err := client.MakeDir(path)
			if err != nil {
				fmt.Println(err)
			}
		case "rmdir":
			err := client.RemoveDirRecur(path)
			if err != nil {
				fmt.Println(err)
			}
		case "cd":
			var err error
			if path == ".." {
				dir = append(dir[:len(dir)-1])
				err = client.ChangeDirToParent()
			} else {
				dir = append(dir, path)
				err = client.ChangeDir(path)
			}
			if err != nil {
				fmt.Println(err)
			}
		case "rmfile":
			err := client.Delete(path)
			if err != nil {
				log.Fatal(err)
			}
		case "ls":
			cur, err := client.CurrentDir()
			if err != nil {
				log.Fatal(err)
			}
			childs, _ := client.List(cur)
			for _, child := range childs {
				fmt.Println(child.Name)
			}
		case "news":
			ti := time.Now().Format(time.RFC1123Z)
			t, _ := time.Parse(time.RFC1123Z, ti)
			fp := gofeed.NewParser()
			feed, _ := fp.ParseURL("http://www.rssboard.org/files/sample-rss-2.xml")
			data := ""
			for _, item := range feed.Items {
				data += item.Title + "\n"
				data += item.Link + "\n"
				data += item.Description + "\n"
				data += item.PublishedParsed.String() + "\n"
				data += item.GUID + "\n\n"
			}
			file, err := os.Create("news.txt")
			if err != nil {
				panic(err)
			}
			file.Write([]byte(data))
			file.Close()

			exist := false
			cur, err := client.CurrentDir() //проверка на уникальность
			if err != nil {
				log.Fatal(err)
			}
			childs, _ := client.List(cur)
			for _, child := range childs {
				if child.Type.String() == "folder" || child.Type.String() == "directory" {
					continue
				}
				serverfile, err := client.Retr(child.Name)
				if err != nil {
					panic(err)
				}
				existfile, err := os.Create(child.Name)
				if err != nil {
					log.Fatal(err)
				}
				_, err = io.Copy(existfile, serverfile)
				if err != nil {
					panic(err)
				}
				serverfile.Close()
				defer serverfile.Close()
				existfile.Close()
				defer existfile.Close()

				existdata, err := os.ReadFile(child.Name)
				if err != nil {
					panic(err)
				}
				err = os.Remove(child.Name)
				if err != nil {
					panic(err)
				}
				if string(existdata) == data {
					fmt.Println("news already loaded")
					exist = true
					break
				}
			}
			if !exist {
				file, err = os.Open("news.txt")
				if err != nil {
					panic(err)
				}
				client.Stor(path+"Донченко_Всеволод_"+t.String()[:10]+"_"+t.String()[11:16]+".txt", file)
			}
			file.Close()
			err = os.Remove("news.txt")
			defer file.Close()
		default:
			fmt.Printf("error: unknown command\n")
			continue

		}
		if exit {
			break
		}
	}
}
func main() {
	//client, err := ftp.Dial("students.yss.su:21", ftp.DialWithTimeout(5*time.Second))
	client, err := ftp.Dial("127.0.0.1:2222", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Login("user", "1234")
	if err != nil {
		log.Fatal(err)
	}
	interact(client)
}

func split(str string) (string, string) {
	if len(str) > 0 {
		if strings.Contains(str, " ") {
			return strings.Split(str, " ")[0], strings.Split(str, " ")[1]
		} else {
			return str, ""
		}
	}
	return "", ""
}
