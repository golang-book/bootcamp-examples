package racey

import (
	"fmt"
	"sync"
)

func F() {
	// c := make(chan int)
	// go func() {
	// 	c <- 5
	// }()
	// fmt.Println(<-c)

	var mu sync.Mutex
	var x int
	go func() {
		mu.Lock()
		x = 5
		mu.Unlock()
	}()
	mu.Lock()
	fmt.Println(x)
	mu.Unlock()
}
