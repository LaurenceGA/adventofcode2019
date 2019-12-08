package main

import (
	"fmt"
	"math"
	"time"

	day7 "github.com/LaurenceGA/adventofcode2019/7"
)

const debug = false

type parameterMode int

const (
	position  = 0
	immediate = 1
)

type opcode int

const (
	add         opcode = 1
	multiply    opcode = 2
	in          opcode = 3
	out         opcode = 4
	jumpIfTrue  opcode = 5
	jumpIfFalse opcode = 6
	lessThan    opcode = 7
	equals      opcode = 8
	halt        opcode = 99
)

func main() {
	start := time.Now()

	max := getMaxAmplifierSignal(day7.GetInputProgram())

	fmt.Println("Max signial is:", max)

	fmt.Println("Time elapsed:", time.Since(start))
}

func getMaxAmplifierSignal(prog []int) int {
	var max int
	for a := 0; a < 5; a++ {
		for b := 0; b < 5; b++ {
			if b == a {
				continue
			}
			for c := 0; c < 5; c++ {
				if c == a || c == b {
					continue
				}
				for d := 0; d < 5; d++ {
					if d == a || d == b || d == c {
						continue
					}
					for e := 0; e < 5; e++ {
						if e == a || e == b || e == c || e == d {
							continue
						}
						out := runAmplifier(prog, amplifierInput{
							a: a,
							b: b,
							c: c,
							d: d,
							e: e,
						})
						if out > max {
							fmt.Println(a, b, c, d, e)
							fmt.Println(out)
							max = out
						}
					}
				}
			}
		}
	}

	return max
}

type amplifierInput struct {
	a, b, c, d, e int
}

func runAmplifier(prog []int, phaseSettings amplifierInput) int {
	progA := copyProgram(prog)
	progB := copyProgram(prog)
	progC := copyProgram(prog)
	progD := copyProgram(prog)
	progE := copyProgram(prog)

	inA := make(chan int)
	aToB := make(chan int)
	bToC := make(chan int)
	cToD := make(chan int)
	dToE := make(chan int)
	outE := make(chan int)

	// Launch programs
	go runProgram(progA, inA, aToB)
	go runProgram(progB, aToB, bToC)
	go runProgram(progC, bToC, cToD)
	go runProgram(progD, cToD, dToE)
	go runProgram(progE, dToE, outE)

	// Phase settings
	inA <- phaseSettings.a
	aToB <- phaseSettings.b
	bToC <- phaseSettings.c
	cToD <- phaseSettings.d
	dToE <- phaseSettings.e

	// Provide initial input
	inA <- 0
	close(inA)

	return <-outE
}

func copyProgram(in []int) (out []int) {
	out = make([]int, len(in))
	copy(out, in)
	return
}

func runProgram(prog []int, input <-chan int, output chan<- int) {
	instructionPointer := 0
	defer close(output)

	if debug {
		fmt.Printf("INITIAL:\n%v IP: %d\n", prog, instructionPointer)
	}

	for {
		instruction := prog[instructionPointer]
		operation := getOpcode(instruction)
		switch operation {
		case halt:
			if debug {
				fmt.Println("HALT")
				fmt.Println(prog)
			}

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

			if debug {
				fmt.Printf("ADD %d(%s)=%d %d(%s)=%d -> %d \n",
					in1,
					paramModeName(paramModes[0]),
					in1Val,
					in2,
					paramModeName(paramModes[1]),
					in2Val,
					out,
				)
			}

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

			if debug {
				fmt.Printf("MUL %d(%s)=%d %d(%s)=%d -> %d \n",
					in1,
					paramModeName(paramModes[0]),
					in1Val,
					in2,
					paramModeName(paramModes[1]),
					in2Val,
					out,
				)
			}

			prog[out] = in1Val * in2Val
			instructionPointer += 4
		case in:
			var out = prog[instructionPointer+1]

			in := <-input
			prog[out] = in

			if debug {
				fmt.Printf("IN %d -> %d\n", in, out)
			}

			instructionPointer += 2
		case out:
			var param = prog[instructionPointer+1]

			paramModes := parameterModes(instruction, 1)

			outVal := getParameterValue(prog, param, paramModes[0])
			output <- outVal

			if debug {
				fmt.Printf("OUT %d(%s)=%d\n",
					param,
					paramModeName(paramModes[0]),
					outVal,
				)
			}

			instructionPointer += 2
		case jumpIfTrue:
			var (
				param           = prog[instructionPointer+1]
				newPointerParam = prog[instructionPointer+2]
			)

			paramModes := parameterModes(instruction, 2)

			var paramVal = getParameterValue(prog, param, paramModes[0])
			var newPointerVal = getParameterValue(prog, newPointerParam, paramModes[1])

			if debug {
				fmt.Printf("JIT %d(%s)=%d -> %d(%s)=%d\n",
					param,
					paramModeName(paramModes[0]),
					paramVal,
					newPointerParam,
					paramModeName(paramModes[1]),
					newPointerVal,
				)
			}

			if paramVal != 0 {
				instructionPointer = newPointerVal
			} else {
				instructionPointer += 3
			}
		case jumpIfFalse:
			var (
				param           = prog[instructionPointer+1]
				newPointerParam = prog[instructionPointer+2]
			)

			paramModes := parameterModes(instruction, 2)

			var paramVal = getParameterValue(prog, param, paramModes[0])
			var newPointerVal = getParameterValue(prog, newPointerParam, paramModes[1])

			if debug {
				fmt.Printf("JIF %d(%s)=%d -> %d(%s)=%d\n",
					param,
					paramModeName(paramModes[0]),
					paramVal,
					newPointerParam,
					paramModeName(paramModes[1]),
					newPointerVal,
				)
			}

			if paramVal == 0 {
				instructionPointer = newPointerVal
			} else {
				instructionPointer += 3
			}
		case lessThan:
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

			if debug {
				fmt.Printf("LT %d(%s)=%d %d(%s)=%d -> %d\n",
					in1,
					paramModeName(paramModes[0]),
					in1Val,
					in2,
					paramModeName(paramModes[1]),
					in2Val,
					out,
				)
			}

			if in1Val < in2Val {
				prog[out] = 1
			} else {
				prog[out] = 0
			}

			instructionPointer += 4
		case equals:
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

			if debug {
				fmt.Printf("EQ %d(%s)=%d %d(%s)=%d -> %d\n",
					in1,
					paramModeName(paramModes[0]),
					in1Val,
					in2,
					paramModeName(paramModes[1]),
					in2Val,
					out,
				)
			}

			if in1Val == in2Val {
				prog[out] = 1
			} else {
				prog[out] = 0
			}

			instructionPointer += 4
		default:
			panic(fmt.Sprintf("Unrecognised opcode %v", operation))
		}
		if debug {
			fmt.Printf("%v IP: %d\n", prog, instructionPointer)
		}
	}
}

func paramModeName(mode parameterMode) string {
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
