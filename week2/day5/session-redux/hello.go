package main

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

func init() {
	//router := http.NewServeMux()
	//router.HandleFunc("/", handleIndex)
	router := httprouter.New()
	router.GET("/:page", handleIndex)
	//router.GET("/users/:name")

	http.Handle("/", context.ClearHandler(router))
	//http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))
}

func handleIndex(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// if I go to: localhost:8080/admin
	// page == "admin"
	params.ByName("page")
}
