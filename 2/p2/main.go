package main

import (
	"fmt"

	day2 "github.com/LaurenceGA/adventofcode2019/2"
)

const targetOutput = 19690720

type opcode int

const (
	add      opcode = 1
	multiply opcode = 2
	halt     opcode = 99
)

func main() {
	for noun := 0; noun < 100; noun++ {
		for verb := 0; verb < 100; verb++ {
			fmt.Printf("noun=%d verb=%d\n", noun, verb)
			program := day2.GetInputProgram(noun, verb)
			runProgram(program)
			output := program[0]

			fmt.Printf("output=%d\n", output)

			if output == targetOutput {
				fmt.Printf("Target goal (%d) reached with noun=%d and verb=%d\n", targetOutput, noun, verb)
				fmt.Println("100*noun+verb =", 100*noun + verb)
				return
			}
		}
	}
}

func runProgram(prog []int) {
	instructionPointer := 0

	for {
		operation := opcode(prog[instructionPointer])
		switch operation {
		case halt:
			return
		case add:
			var (
				in1 = prog[instructionPointer+1]
				in2 = prog[instructionPointer+2]
				out = prog[instructionPointer+3]
			)

			prog[out] = prog[in1] + prog[in2]
			instructionPointer += 4
		case multiply:
			var (
				in1 = prog[instructionPointer+1]
				in2 = prog[instructionPointer+2]
				out = prog[instructionPointer+3]
			)

			prog[out] = prog[in1] * prog[in2]
			instructionPointer += 4
		}
	}
}
