package main

import (
	"fmt"
	"strings"
	"time"

	day6 "github.com/LaurenceGA/adventofcode2019/6"
)

const centreOfMass = "COM"

func main() {
	start := time.Now()

	localOrbits := day6.GetInput()
	fmt.Println(localOrbits)

	com := buildOrbitTree(localOrbits)
	fmt.Println(com.sumHeights(0))

	fmt.Println("time elapsed: ", time.Since(start))
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
