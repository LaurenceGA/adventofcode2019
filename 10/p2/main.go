package main

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"time"

	day10 "github.com/LaurenceGA/adventofcode2019/10"
)

const (
	asteroid = '#'
	empty    = '.'
)

type coord struct {
	X, Y int
}

var station = coord{
	X: 22,
	Y: 19,
}

type ByAngle []coord

func (c ByAngle) Len() int           { return len(c) }
func (c ByAngle) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c ByAngle) Less(i, j int) bool { return angleToStation(c[i]) < angleToStation(c[j]) }

func angleToStation(c coord) float64 {
	straightUpPoint := coord{X: station.X, Y: station.Y - 1}
	p12 := distance(station, straightUpPoint)
	p13 := distance(station, c)
	p23 := distance(straightUpPoint, c)
	angle := math.Acos((math.Pow(p12, 2) + math.Pow(p13, 2) - math.Pow(p23, 2)) / (2 * p12 * p13))
	if c.X < station.X {
		angle = 2*math.Pi - angle
	}
	return angle
}

func main() {
	start := time.Now()

	asteroidMap := day10.GetInput()
	detectable := detectableAsteroids(asteroidMap, station)
	sort.Sort(ByAngle(detectable))

	fmt.Println(detectable)
	goalCoord := detectable[199]
	fmt.Println(goalCoord)
	fmt.Println(goalCoord.X*100 + goalCoord.Y)

	fmt.Println("Time elapsed:", time.Since(start))
}

func detectableAsteroids(asteroidMap [][]rune, pos coord) []coord {
	//printMap(asteroidMap)
	others := otherAsteroids(asteroidMap, pos)
	var detectableAsteroids []coord

	for _, o := range others {
		if canSee(pos, o, asteroidMap) {
			detectableAsteroids = append(detectableAsteroids, o)
		}
	}

	//fmt.Println()
	//printMap(generateDetectionMap(asteroidMap, detectableAsteroids, pos))

	return detectableAsteroids
}

func canSee(a, b coord, asteroidMap [][]rune) bool {
	between := coordsBetween(a, b)
	for _, possibleInterceptor := range between {
		if possibleInterceptor != b && asteroidMap[possibleInterceptor.Y][possibleInterceptor.X] == asteroid {
			return false
		}
	}

	return true
}

func coordsBetween(a, b coord) []coord {
	if a == b {
		return []coord{}
	}

	var intermediateCoords []coord
	xDir, yDir := 1, 1
	if b.X < a.X {
		xDir = -1
	}
	if b.Y < a.Y {
		yDir = -1
	}
	for i := 0; i < int(math.Abs(float64(b.X-a.X)))+1; i++ {
		xCoord := a.X + i*xDir
		for j := 0; j < int(math.Abs(float64(b.Y-a.Y)))+1; j++ {
			yCoord := a.Y + j*yDir
			if (xCoord == a.X && yCoord == a.Y) || (xCoord == b.X && yCoord == b.Y) {
				continue
			}

			curCoord := coord{X: xCoord, Y: yCoord}
			if liesBetween(a, b, curCoord) {
				intermediateCoords = append(intermediateCoords, curCoord)
			}
		}
	}

	return intermediateCoords
}

// Does c lie on the line between a and b?
func liesBetween(a, b, c coord) bool {
	threshold := 0.000001
	distSumDiff := (distance(a, c) + distance(c, b)) - distance(a, b)
	return -threshold < distSumDiff && distSumDiff < threshold
}

func distance(a, b coord) float64 {
	return math.Sqrt(math.Pow(float64(a.X-b.X), 2) + math.Pow(float64(a.Y-b.Y), 2))
}

func getLine(a, b coord) (slope, intercept float64, err error) {
	if (b.X - a.X) == 0 {
		err = errors.New("div by 0")
		return
	}
	slope = float64(b.Y-a.Y) / float64(b.X-a.X)
	intercept = float64(a.Y) - slope*float64(a.X)
	return
}

func otherAsteroids(asteroidMap [][]rune, pos coord) []coord {
	var asteroids []coord
	for y, row := range asteroidMap {
		for x, ast := range row {
			if ast == asteroid && !(x == pos.X && y == pos.Y) {
				asteroids = append(asteroids, coord{X: x, Y: y})
			}
		}
	}

	return asteroids
}

func printMap(asteroidMap [][]rune) {
	for _, row := range asteroidMap {
		for _, ast := range row {
			fmt.Print(string(ast))
		}
		fmt.Println()
	}
}

func generateDetectionMap(originalMap [][]rune, detectables []coord, pos coord) [][]rune {
	newMap := make([][]rune, len(originalMap))
	for i := range newMap {
		newMap[i] = make([]rune, len(originalMap[i]))
		for j := range newMap[i] {
			newMap[i][j] = empty
		}
	}
	for _, obj := range detectables {
		newMap[obj.Y][obj.X] = asteroid
	}

	newMap[pos.Y][pos.X] = 'X'

	return newMap
}
