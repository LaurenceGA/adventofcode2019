package day4

import (
	"strconv"
	"strings"
)

const input = `357253-892942`

// GetInputRange returns the input range
func GetInputRange() (int, int) {
	rangeNums := strings.Split(input, "-")
	lower, err := strconv.Atoi(rangeNums[0])
	if err != nil {
		panic(err)
	}

	upper, err := strconv.Atoi(rangeNums[1])
	if err != nil {
		panic(err)
	}

	return lower, upper
}
