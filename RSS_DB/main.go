package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mmcdole/gofeed"
)

const (
	password string = "Je2dTYr6"
	login    string = "iu9networkslabs"
	host     string = "students.yss.su"
	dbname   string = "iu9networkslabs"
)

func main() {
	db, err := sql.Open("mysql", login+":"+password+"@tcp("+host+")/"+dbname)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("http://www.rssboard.org/files/sample-rss-2.xml")
	for _, item := range feed.Items {
		_, err = db.Exec("insert into iu9networkslabs.iu9Vsevolond (title, link, description, pubDate, guid) values (?, ?, ?, ?, ?)",
			item.Title, item.Link, item.Description, item.PublishedParsed, item.GUID)
		_, err = db.Exec("ALTER TABLE iu9networkslabs.iu9Vsevolond ADD UNIQUE (title, link, description, pubDate, guid)")
		if err != nil {
			panic(err)
		}
	}
}
