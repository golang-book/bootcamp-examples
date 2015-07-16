package example

func Example() {

}

func Add(x, y int) int {
	// xs := []interface{}{
	// 	1, 2, 3, 4, 5,
	// }
	// total := 0
	// for _, x := range xs {
	// 	total += x.(int)
	// }
	// return total
	xs := []int{
		1, 2, 3, 4, 5,
	}
	total := 0
	for _, x := range xs {
		total += x
	}
	return total
}
