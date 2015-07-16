package main

import "fmt"
import "os"
import "strconv"
import "log"

const (
	miTokm  = 1.60934
)

func main() {
  number, err := strconv.ParseFloat(os.Args[1], 64)
  if err != nil {
    log.Fatalln(err)
  }

  fmt.Println("<!DOCTYPE html>")
  fmt.Println("<html>")
  fmt.Println(" <head></head>")
  fmt.Println(" <body>")

  fmt.Printf("Miles: %f<br>\n", number)
  fmt.Printf("Kilometers!: %f<br>\n", number * miTokm)

  fmt.Println(" </body>")
  fmt.Println("</html>")
}
