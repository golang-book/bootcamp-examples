package main

import (
	"fmt"
	"time"
)

func main() {
	timeAsString := "07/10/2015"
	timeAsTime, _ := time.Parse("01/02/2006", timeAsString)

	fmt.Println(timeAsTime.Format(time.UnixDate))

}
