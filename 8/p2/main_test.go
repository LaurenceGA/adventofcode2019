package main

import (
	"fmt"
	"testing"

	day8 "github.com/LaurenceGA/adventofcode2019/8"
)

func TestProcessInput(t *testing.T) {
	in := []byte("123456789012")
	width, height := 3, 2

	input := day8.ProcessInput(in, width, height)
	fmt.Println(input)

	if string(input[0][0]) != "123" {
		t.Errorf("Expected (0,0) to be %s, but got %s", "123", input[0][0])
	}
	if string(input[0][1]) != "456" {
		t.Errorf("Expected (0,1) to be %s, but got %s", "456", input[0][0])
	}
	if string(input[1][0]) != "789" {
		t.Errorf("Expected (1,0) to be %s, but got %s", "789", input[0][0])
	}
	if string(input[1][1]) != "012" {
		t.Errorf("Expected (1,1) to be %s, but got %s", "012", input[0][0])
	}
}

func TestChecksum(t *testing.T) {
	cases := []struct {
		name          string
		img           string
		width, height int
		checksum      int
	}{
		{
			name:     "simple",
			img:      "123456789012",
			width:    3,
			height:   2,
			checksum: 2,
		},
		{
			name:     "more",
			img:      "0210000000102211",
			width:    2,
			height:   2,
			checksum: 4,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			img := day8.ProcessInput([]byte(tt.img), tt.width, tt.height)
			c := checksum(img)
			if c != tt.checksum {
				t.Errorf("Expected: %d, got %d", tt.checksum, c)
			}
		})
	}
}

func TestNumPixels(t *testing.T) {
	cases := []struct {
		name           string
		layer          [][]day8.Pixel
		targetPixel    day8.Pixel
		expectedNumber int
	}{
		{
			name: "2 zeros",
			layer: [][]day8.Pixel{
				{
					day8.Pixel('0'),
					day8.Pixel('3'),
				},
				{
					day8.Pixel('0'),
					day8.Pixel('2'),
				},
			},
			targetPixel:    day8.Pixel('0'),
			expectedNumber: 2,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			num := numPixels(tt.layer, tt.targetPixel)
			if num != tt.expectedNumber {
				t.Errorf("Expected: %d, got %d", tt.expectedNumber, num)
			}
		})
	}
}
