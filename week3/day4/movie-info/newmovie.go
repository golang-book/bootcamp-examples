package movieinfo

import "net/http"

type newMovieModel struct {
}

func handleNewMovie(res http.ResponseWriter, req *http.Request) {
	model := &newMovieModel{}
	err := tpl.ExecuteTemplate(res, "new-movie", model)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
}
