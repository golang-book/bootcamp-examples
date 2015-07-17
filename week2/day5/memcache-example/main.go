package main

import (
	"fmt"
	"net/http"
)
import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/memcache"
)

func init() {

	http.HandleFunc("/", handleIndex)

}

func handleIndex(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	// memcache.Set(ctx, &memcache.Item{
	// 	Key:        "some-key",
	// 	Value:      []byte("some-value"),
	// 	Expiration: 10 * time.Second,
	// })
	item, _ := memcache.Get(ctx, "some-key")
	if item != nil {
		fmt.Fprintln(res, string(item.Value))
	}
	// value := myMap["some-key"]
}
