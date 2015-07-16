package main

import (
	"io"
	"os"
)

func main() {
	src, err := os.Open("src-file")
	if err != nil {
		panic(err)
	}
	defer src.Close()

	dst1, err := os.Create("dst1-file")
	if err != nil {
		panic(err)
	}
	defer dst1.Close()

	if false {
		bs := make([]byte, 5)
		io.ReadFull(src, bs)
		dst1.Write(bs)
	}

	rdr := io.LimitReader(src, 100)
	io.Copy(dst1, rdr)

}
