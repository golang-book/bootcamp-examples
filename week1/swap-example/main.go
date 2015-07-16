package main

import "fmt"

func rotate(args ...*int) {
	if len(args) == 0 {
		return
	}

	firstValue := *args[0]
	for i := 0; i < len(args)-1; i++ {
		*args[i] = *args[i+1]
	}
	*args[len(args)-1] = firstValue
}

var swap = rotate

func main() {
	x := 1
	y := 2
	z := 3
	rotate(&x, &y, &z)
	fmt.Println(x, y, z)
}
