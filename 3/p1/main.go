package main

import (
	"fmt"

	day3 "github.com/LaurenceGA/adventofcode2019/3"
)

func main() {
	wireA, wireB := day3.GetWires()

	wireAPoints := wireA.Rasterize()
	wireBPoints := wireB.Rasterize()

	intersections := findIntersections(wireAPoints, wireBPoints)

	fmt.Println(intersections)
	minDist := 10000
	for _, i := range intersections {
		dist := i.ManhattanDistance()
		if dist < minDist {
			minDist = dist
		}
	}

	fmt.Printf("Min dist is %d\n", minDist)
}

func findIntersections(wireACoords, wireBCoords []*day3.Coordinate) []*day3.Coordinate {
	var intersections []*day3.Coordinate
	for _, cA := range wireACoords {
		if (cA.X != 0 && cA.Y != 0) && !hasCoordinate(intersections, cA) && hasCoordinate(wireBCoords, cA) {
			intersections = append(intersections, cA)
		}
	}

	return intersections
}

func hasCoordinate(s []*day3.Coordinate, c *day3.Coordinate) bool {
	for _, coord := range s {
		if coord.X == c.X && coord.Y == c.Y {
			return true
		}
	}

	return false
}
