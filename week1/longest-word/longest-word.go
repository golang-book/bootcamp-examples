package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func LongestWord(rdr io.Reader) string {
	currentLongestWord := ""
	scanner := bufio.NewScanner(rdr)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		wordOrManyWords := scanner.Text()
		wordOrManyWords = strings.Replace(wordOrManyWords, "?", " ", -1)
		wordOrManyWords = strings.Replace(wordOrManyWords, "-", " ", -1)
		wordOrManyWords = strings.Replace(wordOrManyWords, "/", " ", -1)
		for _, word := range strings.Fields(wordOrManyWords) {
			if len(word) > len(currentLongestWord) {
				currentLongestWord = word
			}
		}
	}
	return currentLongestWord
}

func main() {
	srcFile, err := os.Open("moby.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer srcFile.Close()

	longest := LongestWord(srcFile)
	fmt.Println(longest)
}
