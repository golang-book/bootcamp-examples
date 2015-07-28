package oauthexample

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"

	"github.com/nu7hatch/gouuid"
)

func init() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/github-login", handleGithubLogin)
	http.HandleFunc("/oauth2callback", handleOauth2Callback)
}

func handleIndex(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, `<!DOCTYPE html>
<html>
  <head></head>
  <body>
    <a href="/github-login">LOGIN WITH GITHUB</a>
  </body>
</html>`)
}

func handleGithubLogin(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	// get the session
	session := getSession(ctx, req)
	id, _ := uuid.NewV4()

	redirectURI := "http://localhost:8080/oauth2callback"

	values := make(url.Values)
	values.Add("client_id", "0ccd33716940f347065e")
	values.Add("redirect_uri", redirectURI)
	values.Add("scope", "user:email")
	values.Add("state", id.String())

	// save the session
	session.State = id.String()
	putSession(ctx, res, session)

	http.Redirect(res, req, fmt.Sprintf(
		"https://github.com/login/oauth/authorize?%s",
		values.Encode(),
	), 302)
}

func handleOauth2Callback(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	// get the session
	session := getSession(ctx, req)

	state := req.FormValue("state")
	code := req.FormValue("code")

	if state != session.State {
		http.Error(res, "invalid state", 401)
		return
	}

	accessToken, err := getAccessToken(ctx, state, code)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}

	email, err := getEmail(ctx, accessToken)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}

	session.Email = email
	putSession(ctx, res, session)

	fmt.Fprintln(res, email)

}

func getEmail(ctx context.Context, accessToken string) (string, error) {
	client := urlfetch.Client(ctx)
	response, err := client.Get(
		"https://api.github.com/user/emails?access_token=" + accessToken)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var data []struct {
		Email    string
		Verified bool
		Primary  bool
	}
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return "", err
	}
	if len(data) == 0 {
		return "", fmt.Errorf("no email found")
	}
	return data[0].Email, nil
}

func getAccessToken(ctx context.Context, state, code string) (string, error) {
	values := make(url.Values)
	values.Add("client_id", "0ccd33716940f347065e")
	values.Add("client_secret", "4c21ab338de0449ae13019de25629a7b85e08641")
	values.Add("code", code)
	values.Add("state", state)

	client := urlfetch.Client(ctx)
	response, err := client.PostForm("https://github.com/login/oauth/access_token", values)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	bs, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	values, err = url.ParseQuery(string(bs))
	if err != nil {
		return "", err
	}
	return values.Get("access_token"), nil
}
