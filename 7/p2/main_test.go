package main

import (
	"strconv"
	"sync"
	"testing"

	day7 "github.com/LaurenceGA/adventofcode2019/7"
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
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			prog := day7.ProcessInput(tt.prog)

			inputChannel := make(chan int, len(tt.input))
			for _, i := range tt.input {
				inputChannel <- i
			}
			close(inputChannel)
			outputChannel := make(chan int)

			var wg sync.WaitGroup
			wg.Add(1)
			go runProgram(&wg, prog, inputChannel, outputChannel)

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

func TestAmplifier(t *testing.T) {
	cases := []struct {
		name   string
		prog   string
		input  amplifierInput
		output int
	}{
		{
			name: "1",
			prog: "3,15,3,16,1002,16,10,16,1,16,15,15,4,15,99,0,0",
			input: amplifierInput{
				a: 4,
				b: 3,
				c: 2,
				d: 1,
				e: 0,
			},
			output: 43210,
		},
		{
			name: "2",
			prog: "3,23,3,24,1002,24,10,24,1002,23,-1,23,101,5,23,23,1,24,23,23,4,23,99,0,0",
			input: amplifierInput{
				a: 0,
				b: 1,
				c: 2,
				d: 3,
				e: 4,
			},
			output: 54321,
		},
		{
			name: "3",
			prog: "3,31,3,32,1002,32,10,32,1001,31,-2,31,1007,31,0,33,1002,33,7,33,1,33,31,31,1,32,31,31,4,31,99,0,0,0",
			input: amplifierInput{
				a: 1,
				b: 0,
				c: 4,
				d: 3,
				e: 2,
			},
			output: 65210,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			prog := day7.ProcessInput(tt.prog)

			out := runAmplifier(prog, tt.input)

			if out != tt.output {
				t.Errorf("Expected %d, got %d", tt.output, out)
			}
		})
	}
}

func TestAmplifierMaxmimiser(t *testing.T) {
	cases := []struct {
		name  string
		prog  string
		input amplifierInput
		max   int
	}{
		{
			name: "1",
			prog: "3,15,3,16,1002,16,10,16,1,16,15,15,4,15,99,0,0",
			max:  43210,
		},
		{
			name: "2",
			prog: "3,23,3,24,1002,24,10,24,1002,23,-1,23,101,5,23,23,1,24,23,23,4,23,99,0,0",
			max:  54321,
		},
		{
			name: "3",
			prog: "3,31,3,32,1002,32,10,32,1001,31,-2,31,1007,31,0,33,1002,33,7,33,1,33,31,31,1,32,31,31,4,31,99,0,0,0",
			max:  65210,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			prog := day7.ProcessInput(tt.prog)

			out := getMaxAmplifierSignal(prog)

			if out != tt.max {
				t.Errorf("Expected %d, got %d", tt.max, out)
			}
		})
	}
}

func TestAmplifierFeedbackMaxmimiser(t *testing.T) {
	cases := []struct {
		name  string
		prog  string
		input amplifierInput
		max   int
	}{
		{
			name: "1",
			prog: "3,26,1001,26,-4,26,3,27,1002,27,2,27,1,27,26,27,4,27,1001,28,-1,28,1005,28,6,99,0,0,5",
			max:  139629729,
		},
		{
			name: "2",
			prog: "3,52,1001,52,-5,52,3,53,1,52,56,54,1007,54,5,55,1005,55,26,1001,54,-5,54,1105,1,12,1,53,54,53,1008,54,0,55,1001,55,1,55,2,53,55,53,4,53,1001,56,-1,56,1005,56,6,99,0,0,0,0,10",
			max:  18216,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			prog := day7.ProcessInput(tt.prog)

			out := getMaxFeedbackAmplifierSignal(prog)

			if out != tt.max {
				t.Errorf("Expected %d, got %d", tt.max, out)
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
