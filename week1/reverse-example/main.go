package main

import "fmt"

func main() {
	zs := []int{1, 2, 3, 4, 5}
	fmt.Println(reverse(zs))
}

func reverse(xs []int) []int {
	if len(xs) <= 1 {
		return xs
	}
	lastElement := xs[len(xs)-1]    // 5
	rest := reverse(xs[:len(xs)-1]) // 4, 3, 2, 1
	ys := []int{lastElement}
	ys = append(ys, rest...)
	return ys
}
