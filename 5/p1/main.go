package main

import (
	"fmt"
	"math"

	day5 "github.com/LaurenceGA/adventofcode2019/5"
)

type parameterMode int

const (
	position  = 0
	immediate = 1
)

type opcode int

const (
	add      opcode = 1
	multiply opcode = 2
	in       opcode = 3
	out      opcode = 4
	halt     opcode = 99
)

func main() {
	program := day5.GetInputProgram()
	inputChannel := make(chan int, 1)
	inputChannel <- 1
	close(inputChannel)
	outputChannel := make(chan int)
	go runProgram(program, inputChannel, outputChannel)
	for i := range outputChannel {
		fmt.Printf(">>>>%d<<<<\n", i)
	}
}

func runProgram(prog []int, input <-chan int, output chan<- int) {
	instructionPointer := 0
	defer close(output)

	for {
		instruction := prog[instructionPointer]
		operation := getOpcode(instruction)
		switch operation {
		case halt:
			fmt.Println("HALT")
			return
		case add:
			var (
				in1 = prog[instructionPointer+1]
				in2 = prog[instructionPointer+2]
				out = prog[instructionPointer+3]
			)

			paramModes := parameterModes(instruction, 2)

			var (
				in1Val = getParameterValue(prog, in1, paramModes[0])
				in2Val = getParameterValue(prog, in2, paramModes[1])
			)

			fmt.Printf("ADD %d(%s)=%d %d(%s)=%d -> %d \n",
				in1,
				parmaModeName(paramModes[0]),
				in1Val,
				in2,
				parmaModeName(paramModes[1]),
				in2Val,
				out,
			)

			prog[out] = in1Val + in2Val
			instructionPointer += 4
		case multiply:
			var (
				in1 = prog[instructionPointer+1]
				in2 = prog[instructionPointer+2]
				out = prog[instructionPointer+3]
			)

			paramModes := parameterModes(instruction, 2)

			var (
				in1Val = getParameterValue(prog, in1, paramModes[0])
				in2Val = getParameterValue(prog, in2, paramModes[1])
			)

			fmt.Printf("MUL %d(%s)=%d %d(%s)=%d -> %d \n",
				in1,
				parmaModeName(paramModes[0]),
				in1Val,
				in2,
				parmaModeName(paramModes[1]),
				in2Val,
				out,
			)

			prog[out] = in1Val * in2Val
			instructionPointer += 4
		case in:
			var out = prog[instructionPointer+1]

			in := <-input
			prog[out] = in

			fmt.Printf("IN %d -> %d\n", in, out)

			instructionPointer += 2
		case out:
			var param = prog[instructionPointer+1]

			paramModes := parameterModes(instruction, 1)

			outVal := getParameterValue(prog, param, paramModes[0])
			output <- outVal

			fmt.Printf("OUT %d(%s)=%d\n",
				param,
				parmaModeName(paramModes[0]),
				outVal,
			)

			instructionPointer += 2
		default:
			panic("Unrecognised opcode")
		}
	}
}

func parmaModeName(mode parameterMode) string {
	switch mode {
	case position:
		return "position"
	case immediate:
		return "immediate"
	default:
		panic(fmt.Sprintf("Unknown mode: %d", mode))
	}
}

func getParameterValue(prog []int, parameter int, mode parameterMode) int {
	if mode == immediate {
		return parameter
	}

	return prog[parameter]
}

func parameterModes(instruction, params int) []parameterMode {
	parameterModes := make([]parameterMode, 0, params)
	for i := 2; i < params+2; i++ {
		parameterModes = append(parameterModes, parameterMode(digitAt(instruction, i)))
	}

	return parameterModes
}

func getOpcode(instruction int) opcode {
	return opcode(digitAt(instruction, 1)*10 + digitAt(instruction, 0))
}

func digitAt(num, pos int) int {
	return num / int(math.Pow10(pos)) % 10
}
