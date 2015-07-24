package movieinfo

import (
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/search"
	"google.golang.org/appengine/urlfetch"
)

func renderMarkdown(ctx context.Context, text string) (string, error) {
	client := urlfetch.Client(ctx)
	result, err := client.Post(
		"https://api.github.com/markdown/raw",
		"text/plain",
		strings.NewReader(text),
	)
	if err != nil {
		return "", err
	}
	defer result.Body.Close()
	bs, _ := ioutil.ReadAll(result.Body)
	return string(bs), nil
}

type newMovieModel struct {
	CreatedID string
}

func handleNewMovie(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	model := &newMovieModel{}

	if req.Method == "POST" {
		title := req.FormValue("title")
		summary := req.FormValue("summary")

		summary, err := renderMarkdown(ctx, summary)
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}

		index, err := search.Open("movies")
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}

		movie := &Movie{
			Title:   title,
			Summary: search.HTML(summary),
		}

		id, err := index.Put(ctx, "", movie)
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}
		model.CreatedID = id
	}

	err := tpl.ExecuteTemplate(res, "new-movie", model)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
}
