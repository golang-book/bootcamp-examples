package main

import (
	"log"
	"net/http"
)

func main() {
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
