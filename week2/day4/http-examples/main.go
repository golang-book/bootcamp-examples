package main

import (
	"net/http"

	"github.com/alecthomas/template"
)

func handleSomeRoute(res http.ResponseWriter, req *http.Request) {
	// Step #1: Parse template
	tpl, err := template.ParseFiles(
		"assets/templates/index.gohtml",
		"assets/templates/index.gohtml",
	)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}

	type MyModel struct {
		SomeField int
		Values    []int
		Data      string
	}

	// Step #2: Execute template
	tpl.ExecuteTemplate(res, "assets/templates/index.gohtml", MyModel{
		SomeField: 123,
		Values:    []int{1, 2, 3, 4, 5},
		Data:      "Some Value",
	})

}

func main() {

	http.HandleFunc("/some/route", handleSomeRoute)

	http.ListenAndServe(":8080", nil)
}
