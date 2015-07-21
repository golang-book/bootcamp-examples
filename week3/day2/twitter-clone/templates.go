package main

import (
	"bytes"
	"net/http"
	"time"

	"github.com/dustin/go-humanize"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"html/template"
)

func renderTemplate(res http.ResponseWriter, req *http.Request, name string, data interface{}) {
	// parse templates
	tpl := template.New("")
	tpl = tpl.Funcs(template.FuncMap{
		"humanize_time": func(tm time.Time) string {
			return humanize.Time(tm)
		},
	})
	tpl, err := tpl.ParseGlob("templates/*.gohtml")
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
	// execute page
	var buf bytes.Buffer
	err = tpl.ExecuteTemplate(&buf, name, data)
	if err == nil {
		// execute layout
		type Model struct {
			Body     template.HTML
			LoggedIn bool
		}
		model := Model{
			Body: template.HTML(buf.String()),
		}
		log.Infof(appengine.NewContext(req), "COOKIE: %v", req.Cookies())
		if cookie, err := req.Cookie("logged_in"); err == nil {
			model.LoggedIn = cookie.Value == "true"
		}
		err = tpl.ExecuteTemplate(res, "layout", model)
	}
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
}
