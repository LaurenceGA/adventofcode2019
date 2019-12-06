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
		fuelRequired += fuelRequiredForMass(moduleMass)
	}

	fmt.Println(fuelRequiredForMass(1969))
	fmt.Printf("Total fuel required is %d\n", fuelRequired)
}

func fuelRequiredForMass(mass int) int {
	moduleFuelRequired := int(math.Floor(float64(mass)/3)) - 2

	if moduleFuelRequired <= 0 {
		return 0
	}

	return moduleFuelRequired + fuelRequiredForMass(moduleFuelRequired)
}
