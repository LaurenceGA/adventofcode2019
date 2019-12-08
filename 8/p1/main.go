package main

import (
	"fmt"
	"time"

	day8 "github.com/LaurenceGA/adventofcode2019/8"
)

const (
	width  = 25
	height = 6
)

func main() {
	start := time.Now()

	img := day8.GetInput(width, height)
	fmt.Println(img)

	fmt.Println("Checksum:", checksum(img))

	fmt.Println("Time elapsed:", time.Since(start))
}

func checksum(img [][][]day8.Pixel) int {
	layer := layerWithLeast(img, day8.Pixel('0'))
	numOf1 := numPixels(layer, day8.Pixel('1'))
	numOf2 := numPixels(layer, day8.Pixel('2'))

	return numOf1 * numOf2
}

func layerWithLeast(img [][][]day8.Pixel, leastTarget day8.Pixel) [][]day8.Pixel {
	fmt.Println("least:", leastTarget)
	var minLayer [][]day8.Pixel
	var minPixels int
	for _, l := range img {
		zeros := numPixels(l, leastTarget)
		if minLayer == nil || zeros < minPixels {
			minLayer = l
			minPixels = zeros
		}
	}

	fmt.Println("layer:", minLayer)
	fmt.Println("num:", minPixels)

	return minLayer
}

func numPixels(layer [][]day8.Pixel, targetPixel day8.Pixel) int {
	sum := 0
	for _, r := range layer {
		for _, c := range r {
			if c == targetPixel {
				sum++
			}
		}
	}

	return sum
}
