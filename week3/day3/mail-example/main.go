package main

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/user"

	"google.golang.org/appengine/log"
	"google.golang.org/appengine/mail"
)

func init() {
	http.HandleFunc("/sendmail", handleIndex)
}

func handleIndex(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	u := user.Current(ctx)

	msg := &mail.Message{
		Sender:  u.Email,
		To:      []string{"Caleb Doxsey <caleb@doxsey.net>"},
		Subject: "See you tonight",
		Body:    "Don't forget our plans. Hark, 'til later.",
	}
	if err := mail.Send(ctx, msg); err != nil {
		log.Errorf(ctx, "Alas, my user, the email failed to sendeth: %v", err)
	}
}
