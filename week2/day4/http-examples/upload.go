package main

import (
	"io"
	"net/http"
	"os"
)

func uploadPage(res http.ResponseWriter, req *http.Request) {
	// CHECK FOR VALID SESSION

	if req.Method == "POST" {
		// <input type="file" name="file">
		src, hdr, err := req.FormFile("file")
		if err != nil {

			http.Error(res, err.Error(), 500)
			return
		}
		defer src.Close()

		// create a new file
		dst, err := os.Create("/tmp/" + hdr.Filename)
		if err != nil {

			http.Error(res, err.Error(), 500)
			return
		}
		defer dst.Close()

		// copy the uploaded file into the new file
		io.Copy(dst, src)
	}

	// RENDER THE UPLOAD PAGE
}
