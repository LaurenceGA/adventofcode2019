package main

import (
	"strconv"
	"testing"

	day16 "github.com/LaurenceGA/adventofcode2019/16"
)

func TestAlgorithm(t *testing.T) {
	cases := []struct {
		input          string
		phases         int
		firstXDigits   int
		expectedOutput string
	}{
		{
			input:          `12345678`,
			phases:         1,
			expectedOutput: `48226158`,
		},
		{
			input:          `12345678`,
			phases:         2,
			expectedOutput: `34040438`,
		},
		{
			input:          `12345678`,
			phases:         3,
			expectedOutput: `03415518`,
		},
		{
			input:          `12345678`,
			phases:         4,
			expectedOutput: `01029498`,
		},
		{
			input:          `80871224585914546619083218645595`,
			phases:         100,
			firstXDigits:   8,
			expectedOutput: `24176176`,
		},
		{
			input:          `19617804207202209144916044189917`,
			phases:         100,
			firstXDigits:   8,
			expectedOutput: `73745418`,
		},
		{
			input:          `69317163492948606335995924319873`,
			phases:         100,
			firstXDigits:   8,
			expectedOutput: `52432133`,
		},
	}

	for i, tt := range cases {
		tt := tt
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			nums := day16.ProcessInput(tt.input)
			output := runAlgo(nums, tt.phases)
			if tt.firstXDigits != 0 {
				output = output[:tt.firstXDigits]
			}
			outputString := joinNums(output)
			if outputString != tt.expectedOutput {
				t.Errorf("Expected %s, found %s", tt.expectedOutput, outputString)
			}
		})
	}
}

func TestRepeatPattern(t *testing.T) {
	for i, tt := range []struct {
		inputPattern          []int
		repetitionTimes       int
		expectedOutputPattern []int
	}{
		{
			inputPattern:          []int{0, 1, 0, -1},
			repetitionTimes:       1,
			expectedOutputPattern: []int{0, 1, 0, -1},
		},
		{
			inputPattern:          []int{0, 1, 0, -1},
			repetitionTimes:       2,
			expectedOutputPattern: []int{0, 0, 1, 1, 0, 0, -1, -1},
		},
		{
			inputPattern:          []int{0, 1, 0, -1},
			repetitionTimes:       3,
			expectedOutputPattern: []int{0, 0, 0, 1, 1, 1, 0, 0, 0, -1, -1, -1},
		},
		{
			inputPattern:          []int{0, 1, 0, -1},
			repetitionTimes:       4,
			expectedOutputPattern: []int{0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, -1, -1, -1, -1},
		},
		{
			inputPattern:          []int{0, 0, 0},
			repetitionTimes:       4,
			expectedOutputPattern: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actualOutput := repeatPattern(tt.inputPattern, tt.repetitionTimes)
			actualAsString := joinNums(actualOutput)
			expectedAsString := joinNums(tt.expectedOutputPattern)
			if actualAsString != expectedAsString {
				t.Errorf("Expected %s, got %s", expectedAsString, actualAsString)
			}
		})
	}
}
