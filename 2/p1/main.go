package main

import (
	"fmt"

	day2 "github.com/LaurenceGA/adventofcode2019/2"
)

type opcode int

const (
	add      opcode = 1
	multiply opcode = 2
	halt     opcode = 99
)

func main() {
	program := day2.GetInputProgram(12, 2)
	runProgram(program)

	fmt.Printf("Position 0 is %d\n", program[0])
}

func runProgram(prog []int) {
	curPosition := 0

	for {
		operation := opcode(prog[curPosition])
		switch operation {
		case halt:
			return
		case add:
			var (
				in1 = prog[curPosition+1]
				in2 = prog[curPosition+2]
				out = prog[curPosition+3]
			)

			prog[out] = prog[in1] + prog[in2]
			curPosition += 4
		case multiply:
			var (
				in1 = prog[curPosition+1]
				in2 = prog[curPosition+2]
				out = prog[curPosition+3]
			)

			prog[out] = prog[in1] * prog[in2]
			curPosition += 4
		}
	}
}
