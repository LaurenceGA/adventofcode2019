package main

import (
	"strconv"
	"testing"

	day9 "github.com/LaurenceGA/adventofcode2019/9"
)

func TestProgram(t *testing.T) {
	cases := []struct {
		name   string
		prog   string
		input  []int
		output []int
	}{
		{
			name:   "Equals 8, true, pos",
			prog:   "3,9,8,9,10,9,4,9,99,-1,8",
			input:  []int{8},
			output: []int{1},
		},
		{
			name:   "Equals 8, false, pos",
			prog:   "3,9,8,9,10,9,4,9,99,-1,8",
			input:  []int{10},
			output: []int{0},
		},
		{
			name:   "Less than 8, true, pos",
			prog:   "3,9,7,9,10,9,4,9,99,-1,8",
			input:  []int{5},
			output: []int{1},
		},
		{
			name:   "Less than 8, false (equal), pos",
			prog:   "3,9,7,9,10,9,4,9,99,-1,8",
			input:  []int{8},
			output: []int{0},
		},
		{
			name:   "Less than 8, false (greater), pos",
			prog:   "3,9,7,9,10,9,4,9,99,-1,8",
			input:  []int{9},
			output: []int{0},
		},
		{
			name:   "Equals 8, true, imm",
			prog:   "3,3,1108,-1,8,3,4,3,99",
			input:  []int{8},
			output: []int{1},
		},
		{
			name:   "Equals 8, false, imm",
			prog:   "3,3,1108,-1,8,3,4,3,99",
			input:  []int{10},
			output: []int{0},
		},
		{
			name:   "Less than 8, true, imm",
			prog:   "3,3,1107,-1,8,3,4,3,99",
			input:  []int{5},
			output: []int{1},
		},
		{
			name:   "Less than 8, false (equal), imm",
			prog:   "3,3,1107,-1,8,3,4,3,99",
			input:  []int{8},
			output: []int{0},
		},
		{
			name:   "Less than 8, false (greater), imm",
			prog:   "3,3,1107,-1,8,3,4,3,99",
			input:  []int{9},
			output: []int{0},
		},
		{
			name:   "Is zero true, pos",
			prog:   "3,12,6,12,15,1,13,14,13,4,13,99,-1,0,1,9",
			input:  []int{0},
			output: []int{0},
		},
		{
			name:   "Is zero false, pos",
			prog:   "3,12,6,12,15,1,13,14,13,4,13,99,-1,0,1,9",
			input:  []int{10},
			output: []int{1},
		},
		{
			name:   "Is zero true, imm",
			prog:   "3,3,1105,-1,9,1101,0,0,12,4,12,99,1",
			input:  []int{0},
			output: []int{0},
		},
		{
			name:   "Is zero false, imm",
			prog:   "3,3,1105,-1,9,1101,0,0,12,4,12,99,1",
			input:  []int{10},
			output: []int{1},
		},
		{
			name:   "Larger, below 8",
			prog:   "3,21,1008,21,8,20,1005,20,22,107,8,21,20,1006,20,31,1106,0,36,98,0,0,1002,21,125,20,4,20,1105,1,46,104,999,1105,1,46,1101,1000,1,20,4,20,1105,1,46,98,99",
			input:  []int{5},
			output: []int{999},
		},
		{
			name:   "Larger, equal 8",
			prog:   "3,21,1008,21,8,20,1005,20,22,107,8,21,20,1006,20,31,1106,0,36,98,0,0,1002,21,125,20,4,20,1105,1,46,104,999,1105,1,46,1101,1000,1,20,4,20,1105,1,46,98,99",
			input:  []int{8},
			output: []int{1000},
		},
		{
			name:   "Larger, above 8",
			prog:   "3,21,1008,21,8,20,1005,20,22,107,8,21,20,1006,20,31,1106,0,36,98,0,0,1002,21,125,20,4,20,1105,1,46,104,999,1105,1,46,1101,1000,1,20,4,20,1105,1,46,98,99",
			input:  []int{56},
			output: []int{1001},
		},
		{
			name:   "Copy itself",
			prog:   "109,1,204,-1,1001,100,1,100,1008,100,16,101,1006,101,0,99",
			input:  []int{},
			output: []int{109, 1, 204, -1, 1001, 100, 1, 100, 1008, 100, 16, 101, 1006, 101, 0, 99},
		},
		{
			name:   "16-digit num",
			prog:   "1102,34915192,34915192,7,4,7,99,0",
			input:  []int{},
			output: []int{1219070632396864},
		},
		{
			name:   "big num",
			prog:   "104,1125899906842624,99",
			input:  []int{},
			output: []int{1125899906842624},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			prog := day9.ProcessInput(tt.prog)

			inputChannel := make(chan int, len(tt.input))
			for _, i := range tt.input {
				inputChannel <- i
			}
			close(inputChannel)
			outputChannel := make(chan int)

			comp := NewComputer(prog, inputChannel, outputChannel)
			go comp.Run()

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

			if i != len(tt.output) {
				t.Errorf("Not enough output. Expected %d, got %d\n", len(tt.output), i)
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
