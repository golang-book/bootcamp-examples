package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
)

func init() {
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public/"))))
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/logout", handleLogout)
	http.HandleFunc("/api/", handleAPI)
}

func handleAPI(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	u := user.Current(ctx)

	// switch on the api method
	switch strings.SplitN(req.URL.Path, "/", 3)[2] {
	case "follow":
		// get username to follow
		var username string
		err := json.NewDecoder(req.Body).Decode(&username)
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}

		// get profile of current user
		profile, err := getProfileByEmail(ctx, u.Email)
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}

		// follow
		err = followUser(ctx, profile.Username, username)
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}

		json.NewEncoder(res).Encode(true)
	case "tweets":
		switch req.Method {
		case "POST":
			// look for the profile
			profile, err := getProfileByEmail(ctx, u.Email)
			if err != nil {
				http.Error(res, err.Error(), 500)
				return
			}

			// fill in the user tweet
			var tweet Tweet
			err = json.NewDecoder(req.Body).Decode(&tweet)
			if err != nil {
				http.Error(res, err.Error(), 500)
				return
			}
			tweet.Time = time.Now()
			tweet.Username = profile.Username

			// create the tweet
			err = createTweet(ctx, profile.Email, &tweet)
			if err != nil {
				http.Error(res, err.Error(), 500)
				return
			}
			json.NewEncoder(res).Encode(tweet)
		default:
			http.Error(res, "method not allowed", 405)
		}
	default:
		http.NotFound(res, req)
	}
}

func handleIndex(res http.ResponseWriter, req *http.Request) {
	// for anything but "/" treat it like a user profile
	if req.URL.Path != "/" {
		handleUserProfile(res, req)
		return
	}
	ctx := appengine.NewContext(req)
	u := user.Current(ctx)

	// get recent tweets
	var tweets []*Tweet
	var err error
	if u == nil {
		tweets, err = getTweets(ctx)
	} else {
		tweets, err = getHomeTweets(ctx, u.Email)
	}
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}

	type Model struct {
		Tweets []*Tweet
	}
	model := Model{
		Tweets: tweets,
	}

	renderTemplate(res, req, "index", model)
}

func handleLogin(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	u := user.Current(ctx)

	// look for the user's profile
	profile, err := getProfileByEmail(ctx, u.Email)
	// if it exists redirect
	if err == nil && profile.Username != "" {
		http.SetCookie(res, &http.Cookie{Name: "logged_in", Value: "true"})
		http.Redirect(res, req, "/"+profile.Username, 302)
		return
	}

	type Model struct {
		Profile *Profile
		Error   string
	}
	model := Model{
		Profile: &Profile{Email: u.Email},
	}

	// create the profile
	username := req.FormValue("username")
	if username != "" {
		_, err = getProfileByUsername(ctx, username)
		// if the username is already taken
		if err == nil {
			model.Error = "username is not available"
		} else {
			model.Profile.Username = username
			// try to create the profile
			err = createProfile(ctx, model.Profile)
			if err != nil {
				model.Error = err.Error()
			} else {
				// on success redirect to the user's timeline
				waitForProfile(ctx, username)
				http.SetCookie(res, &http.Cookie{Name: "logged_in", Value: "true"})
				http.Redirect(res, req, "/"+username, 302)
				return
			}
		}
	}

	// render the login template
	renderTemplate(res, req, "login", model)
}

func handleLogout(res http.ResponseWriter, req *http.Request) {
	http.SetCookie(res, &http.Cookie{Name: "logged_in", Value: ""})
	http.Redirect(res, req, "/", 302)
}

func handleUserProfile(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	u := user.Current(ctx)

	// get profile of current user
	myProfile, _ := getProfileByEmail(ctx, u.Email)

	parts := strings.SplitN(req.URL.Path, "/", 3)
	// get the page name
	pageName := ""
	if len(parts) > 2 {
		pageName = parts[2]
	}

	// get the username
	username := parts[1]
	// get the profile
	profile, err := getProfileByUsername(ctx, username)
	if err != nil {
		http.Error(res, err.Error(), 404)
		return
	}

	switch pageName {
	case "":
		tweets, err := getUserTweets(ctx, profile.Username)
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}

		// Render the template
		type Model struct {
			Profile     *Profile
			Tweets      []*Tweet
			CanFollow   bool
			IsFollowing bool
		}
		model := Model{
			Profile: profile,
			Tweets:  tweets,
		}
		if myProfile != nil {
			model.CanFollow = true
			model.IsFollowing = myProfile.IsFollowing(username)
		} else {
			model.CanFollow = false
		}
		renderTemplate(res, req, "user-profile-tweets", model)
	case "following":
		type Model struct {
			Profile     *Profile
			Following   []string
			CanFollow   bool
			IsFollowing bool
		}
		model := Model{
			Profile:   profile,
			Following: profile.Following,
		}
		if myProfile != nil {
			model.CanFollow = true
			model.IsFollowing = myProfile.IsFollowing(username)
		} else {
			model.CanFollow = false
		}
		renderTemplate(res, req, "user-profile-following", model)
	case "followers":
		followers, err := getFollowers(ctx, username)
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}

		type Model struct {
			Profile     *Profile
			Followers   []string
			CanFollow   bool
			IsFollowing bool
		}
		model := Model{
			Profile:   profile,
			Followers: followers,
		}
		if myProfile != nil {
			model.CanFollow = true
			model.IsFollowing = myProfile.IsFollowing(username)
		} else {
			model.CanFollow = false
		}
		renderTemplate(res, req, "user-profile-followers", model)
	default:
		http.NotFound(res, req)
	}

}
