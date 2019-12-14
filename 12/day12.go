package day12

import (
	"math"
	"regexp"
	"strconv"
	"strings"
)

const input = `<x=-4, y=-9, z=-3>
<x=-13, y=-11, z=0>
<x=-17, y=-7, z=15>
<x=-16, y=4, z=2>`

var planetRegex = regexp.MustCompile("^<x=(.*), y=(.*), z=(.*)>$")

// GetInput provides the planet positions from actual puzzle input
func GetInput() (Io *Planet, Europa *Planet, Ganymede *Planet, Callisto *Planet) {
	return ProcessInput(input)
}

// ProcessInput provides the planet positions from arbitrary input
func ProcessInput(in string) (io *Planet, europa *Planet, ganymede *Planet, callisto *Planet) {
	points := strings.Split(in, "\n")
	io = &Planet{
		position: toPoint(points[0]),
	}
	europa = &Planet{
		position: toPoint(points[1]),
	}
	ganymede = &Planet{
		position: toPoint(points[2]),
	}
	callisto = &Planet{
		position: toPoint(points[3]),
	}

	return
}

func toPoint(def string) Point {
	matches := planetRegex.FindStringSubmatch(def)

	x, err := strconv.Atoi(matches[1])
	if err != nil {
		panic("Failed to convert num")
	}
	y, err := strconv.Atoi(matches[2])
	if err != nil {
		panic("Failed to convert num")
	}
	z, err := strconv.Atoi(matches[3])
	if err != nil {
		panic("Failed to convert num")
	}

	return Point{
		X: x,
		Y: y,
		Z: z,
	}
}

// Point represents a position in space
type Point struct {
	X, Y, Z int
}

// Velocity is the change in distance over time
type Velocity struct {
	X, Y, Z int
}

// Planet is a planet
type Planet struct {
	position Point
	velocity Velocity
}

// Move updates a planet's position based on its velocity
func (p *Planet) Move() {
	p.MoveX()
	p.MoveY()
	p.MoveZ()
}

// MoveX moves X
func (p *Planet) MoveX() {
	p.position.X += p.velocity.X
}

// MoveY moves Y
func (p *Planet) MoveY() {
	p.position.Y += p.velocity.Y
}

// MoveZ moves Z
func (p *Planet) MoveZ() {
	p.position.Z += p.velocity.Z
}

func (p *Planet) potentialEnergy() int {
	return absInt(p.position.X) + absInt(p.position.Y) + absInt(p.position.Z)
}

func (p *Planet) kineticEnergy() int {
	return absInt(p.velocity.X) + absInt(p.velocity.Y) + absInt(p.velocity.Z)
}

func absInt(x int) int {
	return int(math.Abs(float64(x)))
}

// Energy gives the total energy of a planet
func (p *Planet) Energy() int {
	return p.potentialEnergy() * p.kineticEnergy()
}

// ExertGravityOn pulls another planet closer
func (p *Planet) ExertGravityOn(other *Planet) {
	p.ExertXGravityOn(other)
	p.ExertYGravityOn(other)
	p.ExertZGravityOn(other)
}

// ExertXGravityOn pulls another planet closer on the x axis
func (p *Planet) ExertXGravityOn(other *Planet) {
	if p.position.X > other.position.X {
		other.velocity.X++
	} else if p.position.X < other.position.X {
		other.velocity.X--
	}
}

// ExertYGravityOn pulls another planet closer on the y axis
func (p *Planet) ExertYGravityOn(other *Planet) {
	if p.position.Y > other.position.Y {
		other.velocity.Y++
	} else if p.position.Y < other.position.Y {
		other.velocity.Y--
	}
}

// ExertZGravityOn pulls another planet closer on the z axis
func (p *Planet) ExertZGravityOn(other *Planet) {
	if p.position.Z > other.position.Z {
		other.velocity.Z++
	} else if p.position.Z < other.position.Z {
		other.velocity.Z--
	}
}

// XPos ...
func (p *Planet) XPos() int {
	return p.position.X
}

// YPos ...
func (p *Planet) YPos() int {
	return p.position.Y
}

// ZPos ...
func (p *Planet) ZPos() int {
	return p.position.Z
}

// XVel ...
func (p *Planet) XVel() int {
	return p.velocity.X
}

// YVel ...
func (p *Planet) YVel() int {
	return p.velocity.Y
}

// ZVel ...
func (p *Planet) ZVel() int {
	return p.velocity.Z
}
