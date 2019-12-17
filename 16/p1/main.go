package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	day16 "github.com/LaurenceGA/adventofcode2019/16"
)

var basePattern = []int{0, 1, 0, -1}

func main() {
	start := time.Now()

	nums := day16.GetInput()
	fmt.Println(joinNums(runAlgo(nums, 100)[:8]))

	fmt.Println("Time elapsed:", time.Since(start))
}

func runAlgo(startingNums []int, numPhases int) []int {
	workingList := startingNums
	for phase := 0; phase < numPhases; phase++ {
		workingList = executePhase(workingList)
	}

	return workingList
}

func executePhase(input []int) []int {
	output := make([]int, 0, len(input))
	for i := 0; i < len(input); i++ {
		output = append(output, calculateOutputElement(input, i+1))
	}

	return output
}

func calculateOutputElement(input []int, numElement int) int {
	pattern := repeatPattern(basePattern, numElement)
	element := 0
	for i, num := range input {
		patternElem := pattern[(i+1)%len(pattern)]
		element += num * patternElem
	}

	element = int(math.Abs(float64(element)))
	lastDigit := element % 10

	return lastDigit
}

func repeatPattern(pattern []int, times int) []int {
	if times == 1 {
		return pattern
	}

	newPattern := make([]int, 0, len(pattern)*times)
	for _, n := range pattern {
		for i := 0; i < times; i++ {
			newPattern = append(newPattern, n)
		}
	}

	return newPattern
}

func joinNums(nums []int) string {
	var sb strings.Builder
	for _, n := range nums {
		sb.WriteString(strconv.Itoa(n))
	}

	return sb.String()
}
