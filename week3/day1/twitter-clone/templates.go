package main

import (
	"bytes"
	"net/http"

	"html/template"
)

func renderTemplate(res http.ResponseWriter, name string, data interface{}) {
	// parse templates
	tpl, err := template.ParseGlob("templates/*.gohtml")
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
	// execute page
	var buf bytes.Buffer
	err = tpl.ExecuteTemplate(&buf, name, data)
	if err == nil {
		// execute layout
		type Page struct {
			Body template.HTML
		}
		err = tpl.ExecuteTemplate(res, "layout", Page{
			Body: template.HTML(buf.String()),
		})
	}
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
}
