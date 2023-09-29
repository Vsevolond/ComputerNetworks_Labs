package main

import (
	"html/template"
	"log"
	"net/http"
)

const INDEX_HTML = `
    <!doctype html>
    <html lang="ru">
        <head>
            <meta charset="utf-8"/>
            <p style="font-weight:800; font-size:30px">Криптовалюты c coinmarketcap.com</p>
			<style>
				a:hover {
					text-decoration: underline;
				}
			</style>
        </head>
        <body>
                {{range .}}
                    <a href="https://coinmarketcap.com{{.Ref}}" style="color: #333 !important;font-size: 24px;font-weight: 700;text-decoration: none;line-height: 34px;" target = "_blank">
                    	{{.Name}}
                    </a>
					{{.Volume}}
                    <br/>
                {{end}}
        </body>
    </html>
    `

var indexHtml = template.Must(template.New("index").Parse(INDEX_HTML))

func serveClient(response http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	log.Println("got request", "Method", request.Method, "Path", path)
	if path != "/" && path != "/index.html" {
		log.Println("invalid path", "Path", path)
		response.WriteHeader(http.StatusNotFound)
	} else if err := indexHtml.Execute(response, downloadNews()); err != nil {
		log.Println("HTML creation failed", "error", err)
	} else {
		log.Println("response sent to client successfully")
	}
}

func main() {
	http.HandleFunc("/", serveClient)
	log.Println("starting listener")
	log.Println("listener failed", "error", http.ListenAndServe("127.0.0.1:6060", nil))
}
