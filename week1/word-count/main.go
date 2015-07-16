package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func WordCount(rdr io.Reader) map[string]int {
	counts := map[string]int{}
	scanner := bufio.NewScanner(rdr)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		word := scanner.Text()
		word = strings.ToLower(word)
		counts[word]++
	}
	return counts
}

func main() {
	srcFile, err := os.Open("moby.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer srcFile.Close()

	counts := WordCount(srcFile)
	fmt.Println("Number of whales:", counts["whale"])
}
