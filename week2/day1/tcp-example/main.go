package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		fmt.Println("SOMEONE CONNECTED", conn.RemoteAddr())

		io.WriteString(conn, fmt.Sprint(time.Now()))

		conn.Close()
	}
}
