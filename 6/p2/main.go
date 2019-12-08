package main

import (
	"fmt"
	"strings"
	"time"

	day6 "github.com/LaurenceGA/adventofcode2019/6"
)

const (
	centreOfMass = "COM"
	you          = "YOU"
	santa        = "SAN"
)

func main() {
	start := time.Now()

	localOrbits := day6.GetInput()
	fmt.Println(localOrbits)

	com := buildOrbitTree(localOrbits)
	pathToYou := com.pathTo(you)
	pathToSanta := com.pathTo(santa)
	fmt.Println(pathToYou)
	fmt.Println(pathToSanta)

	fmt.Println(getMinOrbits(pathToYou, pathToSanta))

	fmt.Println("time elapsed: ", time.Since(start))
}

func getMinOrbits(pathA []day6.Orbit, pathB []day6.Orbit) int {
	sameness := pathsSameness(pathA, pathB)

	return len(pathA) + len(pathB) - sameness*2 - 2
}

func pathsSameness(pathA []day6.Orbit, pathB []day6.Orbit) int {
	smallestLength := len(pathA)
	if len(pathB) < smallestLength {
		smallestLength = len(pathB)
	}

	for i := 0; i < smallestLength; i++ {
		if pathA[i] != pathB[i] {
			return i
		}
	}

	return smallestLength
}

func buildOrbitTree(orbits []day6.Orbit) *orbitTree {
	orbitMap := make(map[string]*orbitTree)

	for _, orb := range orbits {
		if _, ok := orbitMap[orb.Orbiting]; !ok {
			orbitMap[orb.Orbiting] = &orbitTree{
				name: orb.Orbiting,
			}
		}

		if _, ok := orbitMap[orb.Orbited]; !ok {
			orbitMap[orb.Orbited] = &orbitTree{
				name: orb.Orbited,
			}
		}

		orbitMap[orb.Orbited].orbitedBy = append(orbitMap[orb.Orbited].orbitedBy, orbitMap[orb.Orbiting])
	}

	fmt.Println(orbitMap)

	return orbitMap[centreOfMass]
}

type orbitTree struct {
	name      string
	orbitedBy []*orbitTree
}

func (ot *orbitTree) String() string {
	names := make([]string, 0, len(ot.orbitedBy))
	for _, o := range ot.orbitedBy {
		names = append(names, o.name)
	}

	return fmt.Sprintf("[%v]", strings.Join(names, ", "))
}

func (ot *orbitTree) sumHeights(preceedingHeight int) int {
	sumHeight := preceedingHeight
	for _, o := range ot.orbitedBy {
		sumHeight += o.sumHeights(preceedingHeight + 1)
	}

	return sumHeight
}

func (ot *orbitTree) pathTo(bodyName string) []day6.Orbit {
	for _, o := range ot.orbitedBy {
		if o.name == bodyName {
			return []day6.Orbit{{
				Orbited:  ot.name,
				Orbiting: o.name,
			}}
		}
	}

	for _, o := range ot.orbitedBy {
		path := o.pathTo(bodyName)
		if path != nil {
			return append([]day6.Orbit{{
				Orbited:  ot.name,
				Orbiting: o.name,
			}}, path...)
		}
	}

	return nil
}
