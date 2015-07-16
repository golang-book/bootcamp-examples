package main

import (
	"io"
	"net/http"
	"strings"
)

type DogHandler int

func (h DogHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "DOG")
}

type CatHandler int

func main() {
	//var dog DogHandler
	//var cat CatHandler

	http.Handle("/assets/",
		http.StripPrefix("/assets",
			http.FileServer(http.Dir("./assets"))))

	http.HandleFunc("/cats/", func(res http.ResponseWriter, req *http.Request) {
		catName := strings.Split(req.URL.Path, "/")[2]
		io.WriteString(res, `<!DOCTYPE html>
  <html>
    <head>
      <link href="/assets/css/main.css" type="text/css" rel="stylesheet">
    </head>
    <body>
      Cat Name: <strong>`+catName+`</strong><br>
      <img src="/assets/home-cat.jpg">
    </body>
  </html>`)
	})

	http.ListenAndServe(":9000", nil)
}
