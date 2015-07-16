package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	srcFile, err := os.Open("main.go")
	if err != nil {
		log.Fatalln(err)
	}
	defer srcFile.Close()

	scanner := bufio.NewScanner(srcFile)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			line = strings.ToUpper(string(line[0])) + line[1:]
		}
		//    fmt.Println(">>>", line)
		//    fmt.Println(strings.ToUpper(line[0:1]) + strings.ToLower(line[1:]))
		fmt.Println(line)
	}
}
