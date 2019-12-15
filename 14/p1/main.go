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

func main() {
	start := time.Now()

	formulae := day14.GetInput()
	fmt.Println(formulae)

	fmt.Println("Time elapsed:", time.Since(start))
}

func amountOfOreToProduceFuel(formulae []day14.Formula) int {
	formulaMap := make(map[string]day14.Formula)
	for _, f := range formulae {
		formulaMap[f.Makes.Name] = f
	}

	return oreToMakeSubstance(formulaMap, fuel, 1)
}

func oreToMakeSubstance(formulaMap map[string]day14.Formula, targetSubstance string, amountNeeded int) int {
	fmt.Printf("Trying to make %d %s\n", amountNeeded, targetSubstance)
	if targetSubstance == ore {
		fmt.Printf("Using %d ORE\n", amountNeeded)
		return amountNeeded
	}

	targetFormula := formulaMap[targetSubstance]
	fmt.Println("Target formula:", targetFormula)
	oreNeededForOneFormula := 0
	for _, f := range targetFormula.Requires {
		o := oreToMakeSubstance(formulaMap, f.Name, f.Amount)
		oreNeededForOneFormula += o
		// fmt.Printf("%d ore needed to make %s component of %s formula\n", o, f.Name, targetSubstance)
	}

	timeToUseFormula := int(math.Ceil(float64(amountNeeded) / float64(targetFormula.Makes.Amount)))
	fmt.Printf("Want %d %s, formula makes %d -> using %d times = %d ORE\n",
		amountNeeded,
		targetSubstance,
		targetFormula.Makes.Amount,
		timeToUseFormula,
		oreNeededForOneFormula * timeToUseFormula,
	)

	return oreNeededForOneFormula * timeToUseFormula
}
