package main

import (
	"strconv"
	"testing"

	day5 "github.com/LaurenceGA/adventofcode2019/5"
)

func TestProgram(t *testing.T) {
	cases := []struct {
		prog   string
		input  []int
		output []int
	}{
		{
			prog:   "1002,4,3,4,33",
			input:  []int{},
			output: []int{},
		},
		{
			prog:   "3,0,4,0,99",
			input:  []int{2},
			output: []int{2},
		},
	}

	for i, tt := range cases {
		tt := tt
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			prog := day5.ProcessInput(tt.prog)

			inputChannel := make(chan int, len(tt.input))
			for _, i := range tt.input {
				inputChannel <- i
			}
			close(inputChannel)
			outputChannel := make(chan int)

			go runProgram(prog, inputChannel, outputChannel)
			
			i := 0
			for o := range outputChannel {
				if i > len(tt.output) {
					t.Errorf("Too many outputs. Expected %d, got >=%d\n", len(tt.output), i)
				}
				if o != tt.output[i] {
					t.Errorf("Expected output %d to be %d, got %d\n", i, tt.output[i], o)
				}
				i++
			}
		})
	}
}

func TestDigitAt(t *testing.T) {
	cases := []struct {
		num, pos, ans int
	}{
		{
			num: 12345678,
			pos: 0,
			ans: 8,
		},
		{
			num: 12345678,
			pos: 1,
			ans: 7,
		},
		{
			num: 12345678,
			pos: 2,
			ans: 6,
		},
		{
			num: 12345678,
			pos: 7,
			ans: 1,
		},
		{
			num: 12345678,
			pos: 8,
			ans: 0,
		},
	}

	for i, tt := range cases {
		tt := tt
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			n := digitAt(tt.num, tt.pos)
			if n != tt.ans {
				t.Errorf("Expected (num=%d, pos=%d) %d, got %d\n", tt.num, tt.pos, tt.ans, n)
			}
		})
	}
}
