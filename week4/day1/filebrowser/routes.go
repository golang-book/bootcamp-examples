package filebrowser

import (
	"html/template"
	"net/http"
	"strings"

	"google.golang.org/appengine"
)

var tpls *template.Template

func init() {
	tpls = template.Must(template.ParseGlob("templates/*.html"))
	http.HandleFunc("/", index)
	http.HandleFunc("/browse/", browse)
}

func index(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(res, req)
		return
	}

	ctx := appengine.NewContext(req)

	// get session
	session := getSession(ctx, req)
	// update session
	if req.Method == "POST" {
		session.Bucket = req.FormValue("bucket")
		session.Credentials = req.FormValue("credentials")
		putSession(ctx, res, session)
		// redirect to browse
		http.Redirect(res, req, "/browse/", 302)
		return
	}

	err := tpls.ExecuteTemplate(res, "index.html", nil)
	if err != nil {
		http.Error(res, err.Error(), 500)
	}
}

type browseModel struct {
	Bucket string
	Folder string
}

func browse(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	session := getSession(ctx, req)

	// if no bucket has been chosen
	if session.Bucket == "" {
		http.Redirect(res, req, "/", 302)
		return
	}

	folder := strings.SplitN(req.URL.Path, "/", 3)[2]

	model := &browseModel{
		Bucket: session.Bucket,
		Folder: folder,
	}

	err := tpls.ExecuteTemplate(res, "browse.html", model)
	if err != nil {
		http.Error(res, err.Error(), 500)
	}
}
