package main

import (
	"errors"
	"fmt"
	"math"
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

func main() {
	start := time.Now()

	asteroidMap := day10.GetInput()
	maxDetectable := 0
	for y, row := range asteroidMap {
		for x, ast := range row {
			if ast == asteroid {
				numDetectable := detectableAsteroids(asteroidMap, coord{
					X: x,
					Y: y,
				})
				if numDetectable > maxDetectable {
					maxDetectable = numDetectable
				}
			}
		}
	}

	fmt.Println(maxDetectable)

	fmt.Println("Time elapsed:", time.Since(start))
}

func detectableAsteroids(asteroidMap [][]rune, pos coord) int {
	others := otherAsteroids(asteroidMap, pos)

	sum := 0
	for _, o := range others {
		if canSee(pos, o, asteroidMap) {
			sum++
		}
	}

	return sum
}

func canSee(a, b coord, asteroidMap [][]rune) bool {
	// fmt.Printf("Can %v see %v?\n", a, b)
	between := coordsBetween(a, b)
	for _, possibleInterceptor := range between {
		if possibleInterceptor != b && asteroidMap[possibleInterceptor.Y][possibleInterceptor.X] == asteroid {
			// fmt.Println("-1")
			return false
		}
	}

	// fmt.Println("+1")
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
	// fmt.Printf("Does %v lie between %v and %v\n", c, a, b)
	slope, intercept, err := getLine(a, b)

	var onLine bool
	if err != nil {
		onLine = c.X == a.X
	} else {
		onLine = math.Abs(float64(c.Y)-slope*float64(c.X)+intercept) < 0.00000000000000000000001
	}

	// fmt.Println(slope, intercept, err, onLine)

	if b.X > a.X {
		if c.X > b.X || c.X < a.X {
			return false
		}
	}

	if b.X < a.X {
		if c.X < b.X || c.X > a.X {
			return false
		}
	}

	if b.Y > a.Y {
		if c.Y > b.Y || c.Y < a.Y {
			return false
		}
	}

	if b.Y < a.Y {
		if c.Y < b.Y || c.Y > a.Y {
			return false
		}
	}

	return onLine
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
			// if x == 11 && y == 16 {
			// 	fmt.Println(x, y, pos, ast, asteroid)
			// 	fmt.Println(row)
			// }
			if ast == asteroid && !(x == pos.X && y == pos.Y) {
				asteroids = append(asteroids, coord{X: x, Y: y})
			}
		}
	}

	return asteroids
}
