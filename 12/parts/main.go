package main

import (
	"fmt"
	"time"

	day12 "github.com/LaurenceGA/adventofcode2019/12"
)

type planetarySystem []*day12.Planet

func main() {
	start := time.Now()

	io, europa, ganymede, callisto := day12.GetInput()
	system := planetarySystem{io, europa, ganymede, callisto}
	fmt.Println(system)

	fmt.Println("Stable steps:", system.stepsUntilStable())

	fmt.Println("Time elasped:", time.Since(start))
}

func (p planetarySystem) totalEnergy() (sum int) {
	for _, planet := range p {
		sum += planet.Energy()
	}

	return
}

func (p planetarySystem) step() {
	p.stepX()
	p.stepY()
	p.stepZ()
}

func (p planetarySystem) stepX() {
	p.applyGravityX()
	p.movePlanetsX()
}

func (p planetarySystem) stepY() {
	p.applyGravityY()
	p.movePlanetsY()
}

func (p planetarySystem) stepZ() {
	p.applyGravityZ()
	p.movePlanetsZ()
}

func (p planetarySystem) applyGravityX() {
	for i, planetA := range p {
		for _, planetB := range p[i:] {
			if planetA == planetB {
				continue
			}
			planetA.ExertXGravityOn(planetB)
			planetB.ExertXGravityOn(planetA)
		}
	}
}

func (p planetarySystem) applyGravityY() {
	for i, planetA := range p {
		for _, planetB := range p[i:] {
			if planetA == planetB {
				continue
			}
			planetA.ExertYGravityOn(planetB)
			planetB.ExertYGravityOn(planetA)
		}
	}
}

func (p planetarySystem) applyGravityZ() {
	for i, planetA := range p {
		for _, planetB := range p[i:] {
			if planetA == planetB {
				continue
			}
			planetA.ExertZGravityOn(planetB)
			planetB.ExertZGravityOn(planetA)
		}
	}
}

func (p planetarySystem) movePlanetsX() {
	for _, planet := range p {
		planet.MoveX()
	}
}

func (p planetarySystem) movePlanetsY() {
	for _, planet := range p {
		planet.MoveY()
	}
}

func (p planetarySystem) movePlanetsZ() {
	for _, planet := range p {
		planet.MoveZ()
	}
}

func (p planetarySystem) String() string {
	s := ""
	for _, plan := range p {
		s += fmt.Sprintf("%v\n", plan)
	}

	return s
}

func (p planetarySystem) stepsUntilStable() int {
	xChan := make(chan int)
	yChan := make(chan int)
	zChan := make(chan int)
	go func() {
		xChan <- p.getXCycle()
	}()
	go func() {
		yChan <- p.getYCycle()
	}()
	go func() {
		zChan <- p.getZCycle()
	}()
	xCycle := <-xChan
	yCycle := <-yChan
	zCycle := <-zChan

	return LCM(xCycle, yCycle, zCycle)
}

// GCD gives greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// LCM gives Least Common Multiple (LCM) via GCD
func LCM(a, b int, integers ...int) int {
	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
}

func (p planetarySystem) copy() planetarySystem {
	pCopy := make([]*day12.Planet, 0, len(p))
	for _, p := range p {
		newPlan := *p
		pCopy = append(pCopy, &newPlan)
	}

	return pCopy
}

func (p planetarySystem) getXCycle() int {
	pCopy := p.copy()

	xCycle := 0
	for {
		p.stepX()
		xCycle++
		if sameX(p, pCopy) {
			return xCycle
		}
	}
}

func (p planetarySystem) getYCycle() int {
	pCopy := p.copy()

	yCycle := 0
	for {
		p.stepY()
		yCycle++
		if sameY(p, pCopy) {
			return yCycle
		}
	}
}

func (p planetarySystem) getZCycle() int {
	pCopy := p.copy()

	zCycle := 0
	for {
		p.stepZ()
		zCycle++
		if sameZ(p, pCopy) {
			return zCycle
		}
	}
}

func sameX(a, b planetarySystem) bool {
	for i, p := range a {
		if p.XPos() != b[i].XPos() || p.XVel() != b[i].XVel() {
			return false
		}
	}

	return true
}

func sameY(a, b planetarySystem) bool {
	for i, p := range a {
		if p.YPos() != b[i].YPos() || p.YVel() != b[i].YVel() {
			return false
		}
	}

	return true
}

func sameZ(a, b planetarySystem) bool {
	for i, p := range a {
		if p.ZPos() != b[i].ZPos() || p.ZVel() != b[i].ZVel() {
			return false
		}
	}

	return true
}
