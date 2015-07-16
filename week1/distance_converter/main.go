package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	feetToKilometers = 0.0003048
	feetToMeters     = 0.3048
	feetToMiles      = 0.000189394

	kilometersToFeet   = 3280.84
	kilometersToMeters = 1000
	kilometersToMiles  = 0.621371

	milesToFeet       = 5280
	milesToKilometers = 1.60934
	milesToMeters     = 1609.34

	metersToFeet       = 3.28084
	metersToKilometers = 0.001
	metersToMiles      = 0.000621371
)

func main() {
	from := os.Args[1]
	to := os.Args[2]

	result := 0.0

	switch {
	case strings.HasSuffix(from, "mi"):
		miles, _ := strconv.ParseFloat(from[:len(from)-2], 64)
		switch to {
		case "km":
			result = miles * milesToKilometers
		case "m":
			result = miles * milesToMeters
		case "ft":
			result = miles * milesToFeet
		case "mi":
			result = miles
		}
	case strings.HasSuffix(from, "km"):
		kilometers, _ := strconv.ParseFloat(from[:len(from)-2], 64)
		switch to {
		case "km":
			result = kilometers
		case "m":
			result = kilometers * kilometersToMeters
		case "ft":
			result = kilometers * kilometersToFeet
		case "mi":
			result = kilometers * kilometersToMiles
		}
	case strings.HasSuffix(from, "m"):
		meters, _ := strconv.ParseFloat(from[:len(from)-1], 64)
		switch to {
		case "km":
			result = meters * metersToKilometers
		case "m":
			result = meters
		case "ft":
			result = meters * metersToFeet
		case "mi":
			result = meters * metersToMiles
		}
	case strings.HasSuffix(from, "ft"):
		feet, _ := strconv.ParseFloat(from[:len(from)-2], 64)
		switch to {
		case "km":
			result = feet * feetToKilometers
		case "m":
			result = feet * feetToMeters
		case "ft":
			result = feet
		case "mi":
			result = feet * feetToMiles
		}
	default:
		log.Fatalln("unknown from type")
	}

	fmt.Printf("%.2f\n", result)
}
