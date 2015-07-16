package main

import "fmt"

func countClumps(xs []int) int {
	clumps := 0
	inClumps := false
	for i := 1; i < len(xs); i++ {
		curr, prev := xs[i], xs[i-1]
		if inClumps {
			if curr != prev {
				inClumps = false
			}
		} else {
			if curr == prev {
				inClumps = true
				clumps++
			}
		}
	}
	return clumps
}

func main() {
	clumps := countClumps([]int{1, 1, 1, 1, 1})
	fmt.Println(clumps) // expect 1
}
