package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

func getGravatarHash(email string) string {
	//  "  Whatever@example.com   "
	email = strings.TrimSpace(email)
	//  "Whatever@example.com"
	email = strings.ToLower(email)
	//  "whatever@example.com"

	h := md5.New()
	io.WriteString(h, email)
	finalBytes := h.Sum(nil)
	finalString := hex.EncodeToString(finalBytes)
	return finalString
}

func main() {
	fmt.Fprint(os.Stderr, "Enter your name:")
	var name string
	fmt.Scanln(&name)

	fmt.Fprint(os.Stderr, "Enter your email:")
	var email string
	fmt.Scanln(&email)
	gravatarHash := getGravatarHash(email)
	fmt.Println(`<!DOCTYPE html>
<html>
  <head>
    <script>
      console.log("HELLO");
    </script>
  </head>
  <body>
    <h1>` + name + `</h1>
    <img src="http://www.gravatar.com/avatar/` + gravatarHash + `?d=identicon">
  </body>
</html>`)
}
