package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func handleConn(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	i := 0
	headers := map[string]string{}
	var url, method string
	for scanner.Scan() {
		ln := scanner.Text()
		fmt.Println(ln)

		if i == 0 {
			fs := strings.Fields(ln)
			method = fs[0]
			url = fs[1]
			fmt.Println("METHOD", method)
			fmt.Println("URL", url)
		} else {
			// in headers now
			if ln == "" {
				break
			}
			fs := strings.SplitN(ln, ": ", 2)
			headers[fs[0]] = fs[1]
		}

		i++
	}

	// parse body
	if method == "POST" || method == "PUT" {
		// fmt.Println("PARSING BODY")
		// amt, _ := strconv.Atoi(headers["Content-Length"])
		// buf := make([]byte, amt)
		// io.ReadFull(conn, buf)
		//
		// fmt.Println("BODY:", string(buf))
	}

	body := `test<strong>test</strong>`

	io.WriteString(conn, "HTTP/1.1 302 OK\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(body))
	fmt.Fprintf(conn, "Content-Type: text/plain\r\n")
	fmt.Fprintf(conn, "Location: http://www.google.com\r\n")
	io.WriteString(conn, "\r\n")
	io.WriteString(conn, body)
}

func main() {
	server, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatalln(err.Error())
		}
		go handleConn(conn)
	}
}
