package main

import (
	"fmt"
	"math"
	"time"

	day17 "github.com/LaurenceGA/adventofcode2019/17"
)

type parameterMode int

const (
	position  = 0
	immediate = 1
	relative  = 2
)

type opcode int

const (
	add                opcode = 1
	multiply           opcode = 2
	in                 opcode = 3
	out                opcode = 4
	jumpIfTrue         opcode = 5
	jumpIfFalse        opcode = 6
	lessThan           opcode = 7
	equals             opcode = 8
	relativeBaseOffset opcode = 9
	halt               opcode = 99
)

func main() {
	start := time.Now()

	program := day17.GetInputProgram()

	in := make(chan int, 1)
	out := make(chan int, 1)
	comp := NewComputer(program, in, out)
	go comp.Run()

	for c := range out {
		fmt.Print(string(c))
	}

	fmt.Println("Time elapsed:", time.Since(start))
}

// Intcode is an int code computer
type Intcode struct {
	prog   []int
	input  <-chan int
	output chan<- int

	instructionPointer int
	relativeBase       int

	debug bool

	done chan struct{}
}

// NewComputer creates a new Intcode computer
func NewComputer(prog []int, input <-chan int, output chan<- int) *Intcode {
	return &Intcode{
		prog:   createProgMemory(prog),
		input:  input,
		output: output,

		debug: false,

		done: make(chan struct{}, 1),
	}
}

func createProgMemory(in []int) (out []int) {
	out = make([]int, len(in)*10)
	copy(out, in)
	return
}

// Run executes the computer's program
func (i *Intcode) Run() {
	defer close(i.output)

	if i.debug {
		fmt.Printf("INITIAL:\n%v IP: %d RB: %d\n", i.prog, i.instructionPointer, i.relativeBase)
	}

	for {
		instruction := i.prog[i.instructionPointer]
		operation := getOpcode(instruction)
		switch operation {
		case halt:
			if i.debug {
				fmt.Println("HALT")
				fmt.Println(i.prog)
			}

			i.done <- struct{}{}
			close(i.done)

			return
		case add:
			i.add(parameterModes(instruction, 3))
		case multiply:
			i.multiply(parameterModes(instruction, 3))
		case in:
			i.in(parameterModes(instruction, 1))
		case out:
			i.out(parameterModes(instruction, 1))
		case jumpIfTrue:
			i.jumpIfTrue(parameterModes(instruction, 2))
		case jumpIfFalse:
			i.jumpIfFalse(parameterModes(instruction, 2))
		case lessThan:
			i.lessThan(parameterModes(instruction, 3))
		case equals:
			i.equals(parameterModes(instruction, 3))
		case relativeBaseOffset:
			i.relativeBaseOffset(parameterModes(instruction, 1))
		default:
			panic(fmt.Sprintf("Unrecognised opcode %v", operation))
		}
		if i.debug {
			fmt.Printf("IP: %d RB: %d\n" /*i.prog,*/, i.instructionPointer, i.relativeBase)
		}
	}
}

func (i *Intcode) add(paramModes []parameterMode) {
	var (
		in1 = i.prog[i.instructionPointer+1]
		in2 = i.prog[i.instructionPointer+2]
		out = i.prog[i.instructionPointer+3]
	)

	var (
		in1Val = i.getParameterValue(in1, paramModes[0])
		in2Val = i.getParameterValue(in2, paramModes[1])
		outVal = i.getAddressParameterValue(out, paramModes[2])
	)

	if i.debug {
		fmt.Printf("ADD %d(%s)=%d %d(%s)=%d -> %d(%s)=%d\n",
			in1,
			paramModeName(paramModes[0]),
			in1Val,
			in2,
			paramModeName(paramModes[1]),
			in2Val,
			out,
			paramModeName(paramModes[2]),
			outVal,
		)
	}

	i.prog[outVal] = in1Val + in2Val
	i.instructionPointer += 4
}

func (i *Intcode) multiply(paramModes []parameterMode) {
	var (
		in1 = i.prog[i.instructionPointer+1]
		in2 = i.prog[i.instructionPointer+2]
		out = i.prog[i.instructionPointer+3]
	)

	var (
		in1Val = i.getParameterValue(in1, paramModes[0])
		in2Val = i.getParameterValue(in2, paramModes[1])
		outVal = i.getAddressParameterValue(out, paramModes[2])
	)

	if i.debug {
		fmt.Printf("MUL %d(%s)=%d %d(%s)=%d -> %d(%s)=%d \n",
			in1,
			paramModeName(paramModes[0]),
			in1Val,
			in2,
			paramModeName(paramModes[1]),
			in2Val,
			out,
			paramModeName(paramModes[2]),
			outVal,
		)
	}

	i.prog[outVal] = in1Val * in2Val
	i.instructionPointer += 4
}

func (i *Intcode) in(paramModes []parameterMode) {
	var out = i.prog[i.instructionPointer+1]
	var outVal = i.getAddressParameterValue(out, paramModes[0])

	in := <-i.input
	i.prog[outVal] = in

	if i.debug {
		fmt.Printf("IN %d -> %d(%s)=%d\n", in, out, paramModeName(paramModes[0]), outVal)
	}

	i.instructionPointer += 2
}

func (i *Intcode) out(paramModes []parameterMode) {
	var param = i.prog[i.instructionPointer+1]

	outVal := i.getParameterValue(param, paramModes[0])
	i.output <- outVal

	if i.debug {
		fmt.Printf("OUT %d(%s)=%d\n",
			param,
			paramModeName(paramModes[0]),
			outVal,
		)
	}

	i.instructionPointer += 2
}

func (i *Intcode) jumpIfTrue(paramModes []parameterMode) {
	var (
		param           = i.prog[i.instructionPointer+1]
		newPointerParam = i.prog[i.instructionPointer+2]
	)

	var paramVal = i.getParameterValue(param, paramModes[0])
	var newPointerVal = i.getParameterValue(newPointerParam, paramModes[1])

	if i.debug {
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
		i.instructionPointer = newPointerVal
	} else {
		i.instructionPointer += 3
	}
}

func (i *Intcode) jumpIfFalse(paramModes []parameterMode) {
	var (
		param           = i.prog[i.instructionPointer+1]
		newPointerParam = i.prog[i.instructionPointer+2]
	)

	var paramVal = i.getParameterValue(param, paramModes[0])
	var newPointerVal = i.getParameterValue(newPointerParam, paramModes[1])

	if i.debug {
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
		i.instructionPointer = newPointerVal
	} else {
		i.instructionPointer += 3
	}
}

func (i *Intcode) lessThan(paramModes []parameterMode) {
	var (
		in1 = i.prog[i.instructionPointer+1]
		in2 = i.prog[i.instructionPointer+2]
		out = i.prog[i.instructionPointer+3]
	)

	var (
		in1Val = i.getParameterValue(in1, paramModes[0])
		in2Val = i.getParameterValue(in2, paramModes[1])
		outVal = i.getAddressParameterValue(out, paramModes[2])
	)

	if i.debug {
		fmt.Printf("LT %d(%s)=%d %d(%s)=%d -> %d(%s)=%d\n",
			in1,
			paramModeName(paramModes[0]),
			in1Val,
			in2,
			paramModeName(paramModes[1]),
			in2Val,
			out,
			paramModeName(paramModes[2]),
			outVal,
		)
	}

	if in1Val < in2Val {
		i.prog[outVal] = 1
	} else {
		i.prog[outVal] = 0
	}

	i.instructionPointer += 4
}

func (i *Intcode) equals(paramModes []parameterMode) {
	var (
		in1 = i.prog[i.instructionPointer+1]
		in2 = i.prog[i.instructionPointer+2]
		out = i.prog[i.instructionPointer+3]
	)

	var (
		in1Val = i.getParameterValue(in1, paramModes[0])
		in2Val = i.getParameterValue(in2, paramModes[1])
		outVal = i.getAddressParameterValue(out, paramModes[2])
	)

	if i.debug {
		fmt.Printf("EQ %d(%s)=%d %d(%s)=%d -> %d(%s)=%d\n",
			in1,
			paramModeName(paramModes[0]),
			in1Val,
			in2,
			paramModeName(paramModes[1]),
			in2Val,
			out,
			paramModeName(paramModes[2]),
			outVal,
		)
	}

	if in1Val == in2Val {
		i.prog[outVal] = 1
	} else {
		i.prog[outVal] = 0
	}

	i.instructionPointer += 4
}

func (i *Intcode) relativeBaseOffset(paramModes []parameterMode) {
	param := i.prog[i.instructionPointer+1]

	paramVal := i.getParameterValue(param, paramModes[0])

	if i.debug {
		fmt.Printf("RBO %d(%s)=%d\n",
			param,
			paramModeName(paramModes[0]),
			paramVal,
		)
	}

	i.relativeBase += paramVal
	i.instructionPointer += 2
}

func (i *Intcode) getParameterValue(parameter int, mode parameterMode) int {
	switch mode {
	case immediate:
		return parameter
	case position:
		return i.prog[parameter]
	case relative:
		return i.prog[i.relativeBase+parameter]
	default:
		panic(fmt.Sprintf("Unknown parameter mode '%d'", mode))
	}
}

func (i *Intcode) getAddressParameterValue(parameter int, mode parameterMode) int {
	switch mode {
	case position:
		fallthrough
	case immediate:
		return parameter
	case relative:
		return i.relativeBase + parameter
	default:
		panic(fmt.Sprintf("Unknown parameter mode '%d'", mode))
	}
}

func paramModeName(mode parameterMode) string {
	switch mode {
	case position:
		return "position"
	case immediate:
		return "immediate"
	case relative:
		return "relative"
	default:
		panic(fmt.Sprintf("Unknown mode: %d", mode))
	}
}

func parameterModes(instruction, params int) []parameterMode {
	parameterModes := make([]parameterMode, 0, params)
	for i := 2; i < params+2; i++ {
		if i == 4 {
			m := parameterMode(digitAt(instruction, i))
			if m == position {
				m = immediate
			}
			parameterModes = append(parameterModes, m)
			continue
		}
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
