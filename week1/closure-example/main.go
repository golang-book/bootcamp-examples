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
	fmt.Println(nextOdd()) // 0
	fmt.Println(nextOdd()) // 2
	fmt.Println(nextOdd()) // 4
}
