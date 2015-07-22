package main

import (
	"fmt"
	"sort"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
)

// Profile is a user's profile
type Profile struct {
	Username  string
	Email     string
	Following []string
}

func (p Profile) IsFollowing(username string) bool {
	for _, f := range p.Following {
		if f == username {
			return true
		}
	}
	return false
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

type TweetsByTime []*Tweet

func (t TweetsByTime) Len() int {
	return len(t)
}
func (t TweetsByTime) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
func (t TweetsByTime) Less(i, j int) bool {
	return t[i].Time.After(t[j].Time)
}

func createTweet(ctx context.Context, email string, tweet *Tweet) error {
	if len(tweet.Text) == 0 {
		return fmt.Errorf("invalid tweet: must contain at least one character")
	}
	if len(tweet.Text) > 140 {
		return fmt.Errorf("invalid tweet: maximum size is 140 characters")
	}
	profileKey := datastore.NewKey(ctx, "Profile", email, 0, nil)
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

func getHomeTweets(ctx context.Context, email string) ([]*Tweet, error) {
	profile, err := getProfileByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	var allTweets []*Tweet
	for _, f := range profile.Following {
		userTweets, err := getUserTweets(ctx, f)
		if err != nil {
			return nil, err
		}
		allTweets = append(allTweets, userTweets...)
	}
	sort.Sort(TweetsByTime(allTweets))
	if len(allTweets) > 10 {
		allTweets = allTweets[:10]
	}
	return allTweets, nil
}

func getUserTweets(ctx context.Context, username string) ([]*Tweet, error) {
	q := datastore.NewQuery("Tweet")
	if username != "" {
		q = q.Filter("Username =", username)
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

func getFollowers(ctx context.Context, username string) ([]string, error) {
	q := datastore.NewQuery("Profile").
		Filter("Following=", username).
		Order("Username")

	followers := []string{}

	it := q.Run(ctx)
	for {
		var profile Profile
		_, err := it.Next(&profile)
		if err == datastore.Done {
			break
		} else if err != nil {
			return nil, err
		}
		followers = append(followers, profile.Username)
	}
	return followers, nil
}

func followUser(ctx context.Context, follower, followee string) error {
	profile, err := getProfileByUsername(ctx, follower)
	if err != nil {
		return err
	}
	// already following just finish
	for _, f := range profile.Following {
		if f == followee {
			return nil
		}
	}
	profile.Following = append(profile.Following, followee)
	return createProfile(ctx, profile)
}
