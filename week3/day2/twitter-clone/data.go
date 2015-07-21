package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
)

// Profile is a user's profile
type Profile struct {
	Username string
	Email    string
}

func createProfile(ctx context.Context, profile *Profile) error {
	key := datastore.NewKey(ctx, "Profile", profile.Email, 0, nil)
	_, err := datastore.Put(ctx, key, profile)
	return err
}

func getProfileByEmail(ctx context.Context, email string) (*Profile, error) {
	key := datastore.NewKey(ctx, "Profile", email, 0, nil)
	var profile Profile
	err := datastore.Get(ctx, key, &profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func getProfileByUsername(ctx context.Context, username string) (*Profile, error) {
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

func waitForProfile(ctx context.Context, username string) error {
	deadline := time.Now().Add(time.Second * 10)
	for time.Now().Before(deadline) {
		_, err := getProfileByUsername(ctx, username)
		if err == nil {
			return nil
		}
		time.Sleep(time.Second * 1)
	}
	return nil
}

// Tweet represents a single, public user tweet
type Tweet struct {
	ID       int64 `datastore:"-"`
	Username string
	Text     string
	Time     time.Time
}

func createTweet(ctx context.Context, tweet *Tweet) error {
	profileKey := datastore.NewKey(ctx, "Profile", tweet.Username, 0, nil)
	key := datastore.NewIncompleteKey(ctx, "Tweet", profileKey)
	key, err := datastore.Put(ctx, key, tweet)
	if err != nil {
		return err
	}
	tweet.ID = key.IntID()
	return nil
}

func getTweets(ctx context.Context) ([]*Tweet, error) {
	return getUserTweets(ctx, "")
}

func getUserTweets(ctx context.Context, username string) ([]*Tweet, error) {
	q := datastore.NewQuery("Tweet")
	if username != "" {
		profileKey := datastore.NewKey(ctx, "Profile", username, 0, nil)
		q = q.Ancestor(profileKey)
	}
	q = q.Order("-Time").Limit(10)
	var tweets []*Tweet
	keys, err := q.GetAll(ctx, &tweets)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(tweets); i++ {
		tweets[i].ID = keys[i].IntID()
	}
	return tweets, nil
}
