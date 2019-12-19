package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
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

type tileID int

const (
	scaffold tileID = 35
	empty    tileID = 46
	newline  tileID = 10
)

type coord struct {
	X, Y int
}

func (c coord) alignmentParam() int {
	return c.X * c.Y
}

type routine int

const (
	A routine = 'A'
	B routine = 'B'
	C routine = 'C'
)

func main() {
	start := time.Now()

	program := day17.GetInputProgram()
	program[0] = 2

	in := make(chan int, 1)
	out := make(chan int, 1)
	comp := NewComputer(program, in, out)
	go comp.Run()

	scaffolding := getScaffolding(out)
	printScaffolding(scaffolding)
	path := getPrimaryPath(scaffolding)
	fmt.Println(path)
	expandedPath := expandPath(path)
	fmt.Println(expandedPath)
	fmt.Println(findFactors(expandedPath))

	fmt.Println("Time elapsed:", time.Since(start))
	return

	// Main:
	fmt.Println(pullString(out, 5))
	//A
	//<R6>
	//<L10>
	//<R8>
	//<R8>
	//<R6>

	//B
	//<6>
	//<L8>
	//<L10>

	//A
	//<R6>
	//<L10>
	//<R8>
	//<R8>
	//<R6>

	// C
	//<6>
	//<L10>
	//<R6>
	//<L10>
	//<R6>

	//B
	//<6>
	//<L8>
	//<L10>

	//--
	//<R12>
	//<L10>
	//<R6>
	//<L10>
	//--

	//<R6>
	//<L10>
	//<R8>
	//<R8>
	//<R6>

	//<6>
	//<L8>
	//<L10>

	//<R6>
	//<L10>
	//<R8>
	//<R8>
	//<R6>

	//<6>
	//<L8>

	//<2,R6>

	//<L8>
	//<2>
	feedInput(in, "A,B,A,C,B")
	in <- int(newline)
	// Function A
	fmt.Println(pullString(out, 12))
	feedInput(in, "R,6,L,10,R")
	in <- int(newline)
	// Function B
	fmt.Println(pullString(out, 12))
	feedInput(in, "8,R,8,R,12,L")
	in <- int(newline)
	// Function C
	fmt.Println(pullString(out, 12))
	feedInput(in, "6,L,10,R,6,L,10,R,6")
	in <- int(newline)

	fmt.Println(pullString(out, 23))
	in <- 'y'
	in <- int(newline)

	go func() {
		for {
			printScreen(out)
		}
	}()
	//fmt.Println(<-out)
	<-comp.done

	fmt.Println("Time elapsed:", time.Since(start))
}

func findFactors(wholeList []string) [][]string {
	wholeStr := strings.Join(wholeList, "")
	maxLen := 60
	for iSize := 19; iSize < maxLen; iSize++ {
		for jSize := 20; jSize < maxLen; jSize++ {
			// fmt.Println(iSize, jSize)
			for kSize := 30; kSize < maxLen; kSize++ {
				fmt.Println(iSize, jSize, kSize)
				// for i := 0; i < len(wholeList)-iSize-jSize-kSize; i++ {
				for j := iSize; j < len(wholeList)-kSize-jSize; j++ {
					for k := j + jSize; k < len(wholeList)-kSize; k++ {
						ws := wholeStr

						iStr := ws[0 : 0+iSize]
						jStr := ws[j : j+jSize]
						kStr := ws[k : k+kSize]
						ws = strings.ReplaceAll(ws, iStr, "A")
						ws = strings.ReplaceAll(ws, jStr, "B")
						ws = strings.ReplaceAll(ws, kStr, "C")
						for _, c := range ws {
							if c != 'A' && c != 'B' && c != 'D' {
								goto fail
							}
						}
						fmt.Println(ws)
						fmt.Println(len(wholeList), len(ws))
						fmt.Println(0, iSize, iStr)
						fmt.Println(j, jSize, jStr)
						fmt.Println(k, kSize, kStr)
					fail:
					}
				}
				// }
			}
		}
	}

	return [][]string{}
}

func sliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i, aVal := range a {
		if aVal != b[i] {
			return false
		}
	}

	return true
}

func expandPath(path string) []string {
	splitPath := strings.Split(path, ",")
	splitPath = splitPath[:len(splitPath)-1]
	fmt.Println(splitPath)
	var newPath []string
	for _, p := range splitPath {
		if p == "L" || p == "R" {
			newPath = append(newPath, p)
		} else {
			num, err := strconv.Atoi(p)
			if err != nil {
				panic(err)
			}
			for i := 0; i < num; i++ {
				newPath = append(newPath, "1")
			}
		}
	}

	return newPath
}

var (
	turnLeft = map[int]int{
		'^': '<',
		'v': '>',
		'>': '^',
		'<': 'v',
	}
	turnRight = map[int]int{
		'^': '>',
		'v': '<',
		'>': 'v',
		'<': '^',
	}
)

func getPrimaryPath(scaffolding [][]int) string {
	var sb strings.Builder
	robotPos, robot := getRobot(scaffolding)
	dist := 0
	for {
		if isScaffold(robotPos.plus(robot, 1), scaffolding) {
			robotPos = robotPos.plus(robot, 1)
			dist++
			continue
		} else {
			if dist != 0 {
				distStr := strconv.Itoa(dist)
				for _, c := range distStr {
					sb.WriteRune(c)
				}
				sb.WriteRune(',')
				dist = 0
			}
		}

		if isScaffold(robotPos.plus(turnLeft[robot], 1), scaffolding) {
			robot = turnLeft[robot]
			sb.WriteRune('L')
			sb.WriteRune(',')
		} else if isScaffold(robotPos.plus(turnRight[robot], 1), scaffolding) {
			robot = turnRight[robot]
			sb.WriteRune('R')
			sb.WriteRune(',')
		} else {
			break
		}
	}

	return sb.String()
}

func isScaffold(c coord, scaffolding [][]int) bool {
	if c.Y >= len(scaffolding) || c.Y < 0 || c.X >= len(scaffolding[c.Y]) || c.X < 0 {
		return false
	}
	if scaffolding[c.Y][c.X] == int(scaffold) {
		return true
	}

	return false
}

func (c coord) plus(robotDir int, dist int) coord {
	switch robotDir {
	case '>':
		return coord{
			X: c.X + dist,
			Y: c.Y,
		}
	case '<':
		return coord{
			X: c.X - dist,
			Y: c.Y,
		}
	case '^':
		return coord{
			X: c.X,
			Y: c.Y - dist,
		}
	case 'v':
		return coord{
			X: c.X,
			Y: c.Y + dist,
		}
	default:
		panic("AAAH")
	}
}

func getRobot(scaffolding [][]int) (coord, int) {
	for y, row := range scaffolding {
		for x, c := range row {
			if c == '>' || c == '<' || c == '^' || c == 'v' {
				return coord{
					X: x,
					Y: y,
				}, c
			}
		}
	}

	panic("AAH NO ROBOT")
}

func feedInput(in chan int, inputStr string) {
	for _, c := range inputStr {
		in <- int(c)
	}
}

func pullString(out chan int, length int) string {
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteString(string(<-out))
	}

	return sb.String()
}

func printScreen(out chan int) {
	scaffolding := getScaffolding(out)

	printScaffolding(scaffolding)
}

func printScaffolding(scaffolding [][]int) {
	for _, row := range scaffolding {
		for _, c := range row {
			fmt.Print(string(c))
		}
		fmt.Println()
	}
}

func getScaffolding(out chan int) [][]int {
	var scaffolding [][]int
	y := 0
	scaffolding = append(scaffolding, []int{})
	for i := 0; i < 3199; i++ {
		c := <-out
		if tileID(c) == newline {
			y++
			scaffolding = append(scaffolding, []int{})
			continue
		}
		scaffolding[y] = append(scaffolding[y], c)
	}
	return scaffolding
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
