package main

import (
	"strconv"
	"testing"

	day10 "github.com/LaurenceGA/adventofcode2019/10"
)

const (
	testMap1 = `.#..#
.....
#####
....#
...##`

	testMap2 = `......#.#.
#..#.#....
..#######.
.#.#.###..
.#..#.....
..#....#.#
#..#....#.
.##.#..###
##...#..#.
.#....####`
	testMap3 = `#.#...#.#.
.###....#.
.#....#...
##.#.#.#.#
....#.#.#.
.##..###.#
..#...##..
..##....##
......#...
.####.###.`
	testMap4 = `.#..#..###
####.###.#
....###.#.
..###.##.#
##.##.#.#.
....###..#
..#.#..#.#
#..#.#.###
.##...##.#
.....#.#..`
	testMap5 = `.#..##.###...#######
##.############..##.
.#.######.########.#
.###.#######.####.#.
#####.##.#.##.###.##
..#####..#.#########
####################
#.####....###.#.#.##
##.#################
#####.##.###..####..
..######..##.#######
####.##.####...##..#
.#####..#.######.###
##...#.##########...
#.##########.#######
.####.#.###.###.#.##
....##.##.###..#####
.#.#.###########.###
#.#.#.#####.####.###
###.##.####.##.#..##`
)

func TestAsteroidsDetectable(t *testing.T) {
	for _, tt := range []struct {
		name                string
		asteroidMapInput    string
		pos                 coord
		detectableAsteroids int
	}{
		{
			name:                "1, (1, 0)",
			asteroidMapInput:    testMap1,
			pos:                 coord{X: 1, Y: 0},
			detectableAsteroids: 7,
		},
		{
			name:                "1, (4, 0)",
			asteroidMapInput:    testMap1,
			pos:                 coord{X: 4, Y: 0},
			detectableAsteroids: 7,
		},
		{
			name:                "1, (0, 2)",
			asteroidMapInput:    testMap1,
			pos:                 coord{X: 0, Y: 2},
			detectableAsteroids: 6,
		},
		{
			name:                "1, (1, 2)",
			asteroidMapInput:    testMap1,
			pos:                 coord{X: 1, Y: 2},
			detectableAsteroids: 7,
		},
		{
			name:                "1, (2, 2)",
			asteroidMapInput:    testMap1,
			pos:                 coord{X: 2, Y: 2},
			detectableAsteroids: 7,
		},
		{
			name:                "1, (3, 2)",
			asteroidMapInput:    testMap1,
			pos:                 coord{X: 3, Y: 2},
			detectableAsteroids: 7,
		},
		{
			name:                "1, (4, 2)",
			asteroidMapInput:    testMap1,
			pos:                 coord{X: 4, Y: 2},
			detectableAsteroids: 5,
		},
		{
			name:                "1, (4, 3)",
			asteroidMapInput:    testMap1,
			pos:                 coord{X: 4, Y: 3},
			detectableAsteroids: 7,
		},
		{
			name:                "1, (3, 4)",
			asteroidMapInput:    testMap1,
			pos:                 coord{X: 3, Y: 4},
			detectableAsteroids: 8,
		},
		{
			name:                "1, (4, 4)",
			asteroidMapInput:    testMap1,
			pos:                 coord{X: 4, Y: 4},
			detectableAsteroids: 7,
		},
		{
			name:                "2, (5, 8)",
			asteroidMapInput:    testMap2,
			pos:                 coord{X: 5, Y: 8},
			detectableAsteroids: 33,
		},
		{
			name:                "3, (1, 2)",
			asteroidMapInput:    testMap3,
			pos:                 coord{X: 1, Y: 2},
			detectableAsteroids: 35,
		},
		{
			name:                "4, (6, 3)",
			asteroidMapInput:    testMap4,
			pos:                 coord{X: 6, Y: 3},
			detectableAsteroids: 41,
		},
		{
			name:                "5, (11, 13)",
			asteroidMapInput:    testMap5,
			pos:                 coord{X: 11, Y: 13},
			detectableAsteroids: 210,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			aMap := day10.ProcessInput(tt.asteroidMapInput)
			numDetectable := detectableAsteroids(aMap, tt.pos)

			if numDetectable != tt.detectableAsteroids {
				t.Errorf("Expected %d detectable asteroids, found %d", tt.detectableAsteroids, numDetectable)
			}
		})

	}
}

func TestGetLine(t *testing.T) {
	for _, tt := range []struct {
		a, b             coord
		err              bool
		slope, intercept float64
	}{
		{
			a:         coord{X: 0, Y: 0},
			b:         coord{X: 5, Y: 0},
			slope:     0,
			intercept: 0,
		},
		{
			a:         coord{X: 3, Y: 7},
			b:         coord{X: 5, Y: 11},
			slope:     2,
			intercept: 1,
		},
		{
			a:   coord{X: 4, Y: 0},
			b:   coord{X: 4, Y: 3},
			err: true,
		},
	} {
		slope, intercept, err := getLine(tt.a, tt.b)
		if slope != tt.slope {
			t.Errorf("Expected slope %f, got %f", tt.slope, slope)
		}
		if intercept != tt.intercept {
			t.Errorf("Expected intercept %f, got %f", tt.intercept, intercept)
		}
		if (err != nil) != tt.err {
			t.Errorf("Expected err %v, got %v", tt.err, err != nil)
		}
	}
}

func TestLiesBetween(t *testing.T) {
	for i, tt := range []struct {
		a, b, c  coord
		expected bool
	}{
		{
			a:        coord{X: 0, Y: 0},
			b:        coord{X: 5, Y: 0},
			c:        coord{X: 3, Y: 0},
			expected: true,
		},
		{
			a:        coord{X: 0, Y: 0},
			b:        coord{X: 5, Y: 0},
			c:        coord{X: 3, Y: 1},
			expected: false,
		},
		{
			a:        coord{X: 3, Y: 4},
			b:        coord{X: 1, Y: 0},
			c:        coord{X: 2, Y: 2},
			expected: true,
		},
		{
			a:        coord{X: 0, Y: 0},
			b:        coord{X: 1, Y: 1},
			c:        coord{X: 2, Y: 2},
			expected: false,
		},
		{
			a:        coord{X: 1, Y: 0},
			b:        coord{X: 1, Y: 1},
			c:        coord{X: 1, Y: 2},
			expected: false,
		},
		{
			a:        coord{X: 1, Y: 1},
			b:        coord{X: 2, Y: 1},
			c:        coord{X: 3, Y: 1},
			expected: false,
		},
		{
			a:        coord{X: 1, Y: 0},
			b:        coord{X: 4, Y: 3},
			c:        coord{X: 3, Y: 2},
			expected: true,
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			between := liesBetween(tt.a, tt.b, tt.c)
			if between != tt.expected {
				t.Errorf("Expected between=%v, got %v", tt.expected, between)
			}
		})
	}
}
