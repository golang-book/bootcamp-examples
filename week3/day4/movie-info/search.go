package movieinfo

import "net/http"

type searchModel struct {
	Query string
}

func handleSearch(res http.ResponseWriter, req *http.Request) {
	query := req.FormValue("q")

	model := &searchModel{
		Query: query,
	}
	err := tpl.ExecuteTemplate(res, "search", model)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
}
