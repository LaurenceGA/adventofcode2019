package main

import (
	"fmt"
	"strconv"
	"time"

	day4 "github.com/LaurenceGA/adventofcode2019/4"
)

func main() {
	start := time.Now()
	lower, upper := day4.GetInputRange()
	sum := 0
	for i := lower; i < upper+1; i++ {
		if meetsPasswordCriteria(i) {
			fmt.Printf("%d meets the criteria\n", i)
			sum++
		}
	}

	fmt.Printf("%d numbers meet the criteria\n", sum)

	fmt.Printf("Execution took %s\n", time.Since(start))
}

func meetsPasswordCriteria(num int) bool {
	if !hasSameAdjacents(num) {
		return false
	}

	if !ascending(num) {
		return false
	}

	return true
}

func ascending(num int) bool {
	asStr := strconv.Itoa(num)
	curVal := asStr[0]
	for _, n := range asStr {
		if int(n) < int(curVal) {
			return false
		}

		curVal = byte(n)
	}

	return true
}

func hasSameAdjacents(num int) bool {
	asStr := strconv.Itoa(num)
	for i := 0; i < len(asStr)-1; i++ {
		if asStr[i] == asStr[i+1] {
			return true
		}
	}

	return false
}
