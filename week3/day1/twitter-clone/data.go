package main

import (
	"fmt"
	"net/http"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// Profile is a user's profile
type Profile struct {
	Username string
	Email    string
}

func createProfile(req *http.Request, profile *Profile) error {
	ctx := appengine.NewContext(req)
	key := datastore.NewKey(ctx, "Profile", profile.Email, 0, nil)
	_, err := datastore.Put(ctx, key, profile)
	return err
}

func getProfileByEmail(req *http.Request, email string) (*Profile, error) {
	ctx := appengine.NewContext(req)
	key := datastore.NewKey(ctx, "Profile", email, 0, nil)
	var profile Profile
	err := datastore.Get(ctx, key, &profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func getProfileByUsername(req *http.Request, username string) (*Profile, error) {
	ctx := appengine.NewContext(req)
	q := datastore.NewQuery("Profile").Filter("Username =", username).Limit(1)
	var profiles []Profile
	_, err := q.GetAll(ctx, &profiles)
	if err != nil {
		return nil, err
	}
	if len(profiles) == 0 {
		return nil, fmt.Errorf("profile not found")
	}
	return &profiles[0], nil
}

func waitForProfile(req *http.Request, username string) error {
	deadline := time.Now().Add(time.Second * 10)
	for time.Now().Before(deadline) {
		_, err := getProfileByUsername(req, username)
		if err == nil {
			return nil
		}
		time.Sleep(time.Second * 1)
	}
	return nil
}
