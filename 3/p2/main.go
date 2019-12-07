package main

import (
	"fmt"
	"time"

	day3 "github.com/LaurenceGA/adventofcode2019/3"
)

func main() {
	start := time.Now()

	wireA, wireB := day3.GetWires()

	wireAPoints := wireA.Rasterize()
	wireBPoints := wireB.Rasterize()
	fmt.Println(wireAPoints)
	fmt.Println(wireBPoints)

	intersections := findIntersections(wireAPoints, wireBPoints)

	fmt.Println(intersections)
	minSteps := 999999
	for _, i := range intersections {
		steps := i.Steps
		if steps < minSteps {
			minSteps = steps
		}
	}

	fmt.Printf("Min steps is %d\n", minSteps)
	fmt.Printf("Execution took %s\n", time.Since(start))
}

func findIntersections(wireACoords, wireBCoords []*day3.Coordinate) []*day3.Coordinate {
	var intersections []*day3.Coordinate
	for _, cA := range wireACoords {
		if (cA.X != 0 && cA.Y != 0) && findMatchingCoordinate(intersections, cA) == nil {
			cB := findMatchingCoordinate(wireBCoords, cA)
			if cB != nil {
				intersections = append(intersections, &day3.Coordinate{
					X:     cA.X,
					Y:     cA.Y,
					Steps: cA.Steps + cB.Steps,
				})
			}
		}
	}

	return intersections
}

func findMatchingCoordinate(s []*day3.Coordinate, c *day3.Coordinate) *day3.Coordinate {
	for _, coord := range s {
		if coord.X == c.X && coord.Y == c.Y {
			return coord
		}
	}

	return nil
}
