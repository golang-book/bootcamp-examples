package main

import (
	"fmt"
	"net/http"

	"github.com/alecthomas/template"
	"github.com/gorilla/sessions"
)

func handleSomeRoute(res http.ResponseWriter, req *http.Request) {
	// Step #1: Parse template
	tpl, err := template.ParseFiles(
		"assets/templates/index.gohtml",
		"assets/templates/login.gohtml",
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
	tpl.ExecuteTemplate(res, "assets/templates/login.gohtml", MyModel{
		SomeField: 123,
		Values:    []int{1, 2, 3, 4, 5},
		Data:      "Some Value",
	})

}

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func loginPage(res http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "session")

	if req.Method == "POST" {
		email := req.FormValue("email")
		password := req.FormValue("password")
		if email == "whatever" && password == "some-password" {
			session.Values["logged_in"] = "YES"
		} else {
			http.Error(res, "invalid credentials", 401)
			return
		}
		// save session
		session.Save(req, res)
		// redirect to main page
		http.Redirect(res, req, "/", 302)
		return
	}

	// RENDER LOGIN FORM
}

func logoutPage(res http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "session")
	delete(session.Values, "logged_in")
	session.Save(req, res)
	http.Redirect(res, req, "/", 302)
}

func handleSomeSessionRoute(res http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "session")
	// Get a Value:
	value := session.Values["whatever-i-call-it"]
	// value is an interface{}
	str, _ := session.Values["whatever-i-call-it"].(string)

	fmt.Println(value, str)

	// Set a Value:
	session.Values["whatever-i-call-it"] = "test"
	// Delete a Value:
	delete(session.Values, "whatever-i-call-it")

	// Save it.
	session.Save(req, res)
}

func main() {

	http.HandleFunc("/some/route", handleSomeRoute)
	http.HandleFunc("/some/session/route", handleSomeSessionRoute)

	// to generate cert and key:
	// go run $(go env GOROOT)/src/crypto/tls/generate_cert.go --host=localhost
	http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", nil)
}
