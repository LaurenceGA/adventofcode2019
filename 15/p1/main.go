package main

import (
	"container/list"
	"fmt"
	"math"
	"time"

	day15 "github.com/LaurenceGA/adventofcode2019/15"
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

	program := day15.GetInputProgram()

	in := make(chan int, 1)
	out := make(chan int, 1)
	comp := NewComputer(program, in, out)
	go comp.Run()

	distToTarget(in, out)

	fmt.Println("Time elapsed:", time.Since(start))
}

type movementDirection int

const (
	north = iota + 1
	south
	west
	east
)

func (m movementDirection) String() string {
	switch m {
	case north:
		return "north"
	case south:
		return "south"
	case west:
		return "west"
	case east:
		return "east"
	default:
		return "unknown"
	}
}

type robotStatus int

const (
	hitWall = iota
	movedToEmptySpace
	movedToTarget
)

func (r robotStatus) String() string {
	switch r {
	case hitWall:
		return "hit wall"
	case movedToEmptySpace:
		return "moved to empty space"
	case movedToTarget:
		return "moved to target!"
	default:
		return "unknown"
	}
}

type coord struct {
	X, Y int
}

func distToTarget(in chan<- int, out <-chan int) {
	curPos := coord{
		X: 0,
		Y: 0,
	}
	directionsToTry := list.New()
	directionsToTry.PushFront(curPos)
	var currentPath []movementDirection
	visited := make(map[coord]struct{})
	visited[curPos] = struct{}{}
	for directionsToTry.Len() != 0 {
		elm := directionsToTry.Front()
		coordToTry := elm.Value.(coord)
		fmt.Println(coordToTry)

		for coordToTry != curPos && !isNeighbour(curPos, coordToTry) {
			fmt.Println("Unwinding...")
			if currentPath == nil {
				panic("Trying to unwind, but curPath is nil")
			}
			movingIn := oppositeDirection(currentPath[0])
			in <- int(movingIn)
			s := <-out
			if s != movedToEmptySpace {
				panic("failed unwinding")
			}
			curPos = coordInDirection(curPos, movingIn)
			currentPath = currentPath[1:]
		}
		if curPos != coordToTry {
			// We are now next to the coordToTry
			actualDir := coordDirection(curPos, coordToTry)
			in <- int(actualDir)
			status := <-out
			currentPath = append([]movementDirection{actualDir}, currentPath...)
			visited[coordToTry] = struct{}{}
			fmt.Printf("Moved %v, status is %v\n", actualDir, status)
			switch status {
			case movedToEmptySpace:
				curPos = coordToTry
			case movedToTarget:
				fmt.Println("!!")
				fmt.Println(len(currentPath))
				return
			case hitWall:
				fmt.Println("Hit a wall...")
			}
		}

		possiblePositions := surroundingPositions(in, out, coordToTry)
		fmt.Println(possiblePositions)
		for _, p := range possiblePositions {
			if _, ok := visited[p]; !ok {
				fmt.Println("Adding", p)
				directionsToTry.PushFront(p)
			}
		}

		directionsToTry.Remove(elm)
	}
}

func coordDirection(from, to coord) movementDirection {
	if to == northCoord(from) {
		return north
	}
	if to == southCoord(from) {
		return south
	}
	if to == westCoord(from) {
		return west
	}
	if to == eastCoord(from) {
		return east
	}

	panic("Don't know coord direction!")
}

func isNeighbour(a, b coord) bool {
	if (b.X == a.X+1 && a.Y == b.Y) ||
		(b.X == a.X-1 && a.Y == b.Y) ||
		(b.Y == a.Y+1 && a.X == b.X) ||
		(b.Y == a.Y-1 && a.X == b.X) {
		return true
	}
	return false
}

func surroundingPositions(in chan<- int, out <-chan int, curPos coord) []coord {
	fmt.Println("Checking directions...")
	var positions []coord
	if dirIsValid(in, out, north) {
		positions = append(positions, northCoord(curPos))
	}
	if dirIsValid(in, out, south) {
		positions = append(positions, southCoord(curPos))
	}
	if dirIsValid(in, out, west) {
		positions = append(positions, westCoord(curPos))
	}
	if dirIsValid(in, out, east) {
		positions = append(positions, eastCoord(curPos))
	}
	fmt.Println("Checked")
	return positions
}

func coordInDirection(start coord, dir movementDirection) coord {
	switch dir {
	case north:
		return northCoord(start)
	case south:
		return southCoord(start)
	case west:
		return westCoord(start)
	case east:
		return eastCoord(start)
	}

	panic("Dunno coord in direction")
}

func northCoord(curPos coord) coord {
	return coord{X: curPos.X, Y: curPos.Y - 1}
}

func southCoord(curPos coord) coord {
	return coord{X: curPos.X, Y: curPos.Y + 1}
}

func westCoord(curPos coord) coord {
	return coord{X: curPos.X - 1, Y: curPos.Y}
}

func eastCoord(curPos coord) coord {
	return coord{X: curPos.X + 1, Y: curPos.Y}
}

func dirIsValid(in chan<- int, out <-chan int, direction movementDirection) bool {
	in <- int(direction)
	status := robotStatus(<-out)
	fmt.Println("Moved", direction, "status", status)
	switch status {
	case hitWall:
		return false
	default:
		// move back
		fmt.Println("Moving back")
		in <- int(oppositeDirection(direction))
		newStatus := <-out
		if newStatus != movedToEmptySpace {
			panic("Fail trying to return to previous spot")
		}
		fmt.Println("Moved", oppositeDirection(direction), "status", robotStatus(newStatus))
		return true
	}
}

func oppositeDirection(direction movementDirection) movementDirection {
	switch direction {
	case north:
		return south
	case south:
		return north
	case west:
		return east
	case east:
		return west
	default:
		panic(fmt.Sprintf("Invalid direction %v", direction))
	}
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
