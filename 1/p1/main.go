package main

import (
	"fmt"
	"math"

	day1 "github.com/LaurenceGA/adventofcode2019/1"
)

func main() {
	inputNumbers := day1.GetInputNumbers()

	fuelRequired := 0
	for _, moduleMass := range inputNumbers {
		fuelRequired += getModuleFuelRequired(moduleMass)
	}

	fmt.Printf("Total fuel required is %d\n", fuelRequired)
}

func getModuleFuelRequired(mass int) int {
	return int(math.Floor(float64(mass)/3)) - 2
}
