package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func printHello(i int) {
	fmt.Println("Hello", i)
	wg.Done()
}

func main() {
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go printHello(i)
	}
	wg.Wait()
}
