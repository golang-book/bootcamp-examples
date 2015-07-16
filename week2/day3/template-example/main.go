package main

import (
	"html/template"
	"log"
	"os"
	"strings"
)

type Page struct {
	Title string
	Body  template.HTML
}

func main() {
	log.SetFlags(0)

	var err error

	tpl := template.New("tpl.gohtml")
	tpl = tpl.Funcs(template.FuncMap{
		"uppercase": func(str string) string {
			return strings.ToUpper(str)
		},
	})
	tpl, err = tpl.ParseFiles("tpl.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	err = tpl.ExecuteTemplate(os.Stdout, "tpl.gohtml", Page{
		Title: "My Title 2",
		Body:  "Hello World <script>alert('hi');</script>",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
