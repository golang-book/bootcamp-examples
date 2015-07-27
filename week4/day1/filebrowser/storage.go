package filebrowser

import (
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"google.golang.org/cloud"
	"google.golang.org/cloud/storage"
)

func getCloudContext(aeCtx context.Context) (context.Context, error) {
	data, err := ioutil.ReadFile("/Users/caleb/Desktop/gcs.json")
	if err != nil {
		return nil, err
	}

	conf, err := google.JWTConfigFromJSON(
		data,
		storage.ScopeFullControl,
	)
	if err != nil {
		return nil, err
	}

	tokenSource := conf.TokenSource(aeCtx)

	hc := &http.Client{
		Transport: &oauth2.Transport{
			Source: tokenSource,
			Base:   &urlfetch.Transport{Context: aeCtx},
		},
	}

	return cloud.NewContext(appengine.AppID(aeCtx), hc), nil
}
