package main

import (
	"fmt"
	"math"
	"time"

	day14 "github.com/LaurenceGA/adventofcode2019/14"
)

const (
	ore  = "ORE"
	fuel = "FUEL"
)

const trillion = 1000000000000

func main() {
	start := time.Now()

	formulae := day14.GetInput()
	fmt.Println(maxFuelGenerated(formulae, trillion))

	fmt.Println("Time elapsed:", time.Since(start))
}

func maxFuelGenerated(formulae []day14.Formula, maxOre int) int {
	low := 0
	high := 100000000

	guesses := 0
	finalGuess := 0
	for {
		guesses++
		middle := (high-low)/2 + low
		fmt.Println(high, low, middle)

		oreForGuess := amountOfOreToProduceFuel(formulae, middle)
		fmt.Println(oreForGuess)
		if oreForGuess > trillion {
			fmt.Println("Too high")
			high = middle
		} else {
			fmt.Println("Too low")
			low = middle
		}

		if guesses > 1000 {
			fmt.Println(middle)
			finalGuess = middle
			break
		}
	}

	return finalGuess
}

func amountOfOreToProduceFuel(formulae []day14.Formula, amountNeeded int) int {
	formulaMap := make(map[string]day14.Formula)
	for _, f := range formulae {
		formulaMap[f.Makes.Name] = f
	}

	//fmt.Println(formulae)

	return oreToMakeSubstance(formulaMap, make(map[string]int), fuel, amountNeeded)
}

func oreToMakeSubstance(formulaMap map[string]day14.Formula, excess map[string]int, targetSubstance string, amountNeeded int) int {
	//fmt.Printf("Calculating ORE to make %d %s\n", amountNeeded, targetSubstance)
	if targetSubstance == ore {
		//fmt.Printf("Using %d ORE\n", amountNeeded)
		return amountNeeded
	}

	if _, ok := excess[targetSubstance]; !ok {
		excess[targetSubstance] = 0
	}

	if excess[targetSubstance] >= amountNeeded {
		//fmt.Println("Enough in excess. Taking", amountNeeded)
		excess[targetSubstance] -= amountNeeded
		return 0
	}

	if excess[targetSubstance] > 0 {
		//fmt.Println("Excess present, using", excess[targetSubstance])
		amountNeeded -= excess[targetSubstance]
		excess[targetSubstance] = 0
		//fmt.Println("New desired", targetSubstance, "is", amountNeeded)
	}

	targetFormula := formulaMap[targetSubstance]
	//fmt.Println("Target formula:", targetFormula)

	timesToUseFormula := int(math.Ceil(float64(amountNeeded) / float64(targetFormula.Makes.Amount)))
	//fmt.Printf("Want %d %s, formula makes %d -> using %d times\n",
	//	amountNeeded,
	//	targetSubstance,
	//	targetFormula.Makes.Amount,
	//	timesToUseFormula,
	//)

	oreNeededForOneFormula := 0
	for _, f := range targetFormula.Requires {
		oreNeededForOneFormula += oreToMakeSubstance(formulaMap, excess, f.Name, f.Amount*timesToUseFormula)
	}

	totalSubstanceBeingMade := timesToUseFormula * targetFormula.Makes.Amount

	excess[targetSubstance] += totalSubstanceBeingMade - amountNeeded
	//fmt.Println("Excess:", excess)

	return oreNeededForOneFormula
}
