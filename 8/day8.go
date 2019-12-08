package day8

import (
	"io/ioutil"
)

const inputFile = "8/input.sif"

// Pixel represents a colour
type Pixel byte

// Colours!
const (
	Black       Pixel = '0'
	White       Pixel = '1'
	Transparent Pixel = '2'
)

func (p Pixel) String() string {
	return string(p)
}

// GetInput processes the real puzzle input
func GetInput(width, height int) [][][]Pixel {
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}

	return ProcessInput(data, width, height)
}

// ProcessInput processes arbitrary puzzle input
func ProcessInput(in []byte, width, height int) [][][]Pixel {
	numLayers := len(in) / (width * height)
	data := make([][][]Pixel, 0, numLayers)
	for l := 0; l < numLayers; l++ {
		layer := make([][]Pixel, 0, height)
		for h := 0; h < height; h++ {
			row := in[l*width*height+width*h : l*width*height+width*h+width]

			pixels := make([]Pixel, 0, len(row))
			for _, p := range row {
				pixels = append(pixels, Pixel(p))
			}

			layer = append(layer, pixels)
		}

		data = append(data, layer)
	}

	return data
}
