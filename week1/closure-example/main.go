package main

import "fmt"

func makeOddGenerator() func() int {
	i := int(1)
	return func() int {
		i += 2
		return i
	}
}

func main() {
	nextOdd := makeOddGenerator()
	fmt.Println(nextOdd()) // i = 0 result 3
	fmt.Println(nextOdd()) // i = 2 result 5
	fmt.Println(nextOdd()) // i = 4 result 7
}
