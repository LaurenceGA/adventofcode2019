package main

import (
	"fmt"
	day18 "github.com/LaurenceGA/adventofcode2019/18"
	"time"
)

func main() {
	start := time.Now()

	grid := day18.GetInput()
	for _, g := range grid {
		fmt.Println(g)
	}

	fmt.Println("Time elapsed:", time.Since(start))
}
