package main

import (
	"io"
	"os"
)
import "log"

func main() {
	// my-cat f1 f2 f3 ...
	// os.Args[0] == my-cat
	// os.Args[1] == f1
	// os.Args[2] == f2
	// os.Args[3] == f3
	for _, v := range os.Args[1:] {
		f, err := os.Open(v)
		if err != nil {
			log.Fatalln("my program broke", err.Error())
		}
		defer f.Close()

		io.Copy(os.Stdout, f)
	}

	// $ my-cat
	// Enter File Name:

	//rdr := strings.NewReader("some string")
}
