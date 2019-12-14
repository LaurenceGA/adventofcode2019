package main

import "testing"

import "strconv"
import day12 "github.com/LaurenceGA/adventofcode2019/12"

func TestOrbit(t *testing.T) {
	for i, tt := range []struct {
		input               string
		numSteps            int
		expectedFinalEnergy int
	}{
		{
			input: `<x=-1, y=0, z=2>
<x=2, y=-10, z=-7>
<x=4, y=-8, z=8>
<x=3, y=5, z=-1>`,
			numSteps:            10,
			expectedFinalEnergy: 179,
		},
		{
			input: `<x=-8, y=-10, z=0>
<x=5, y=5, z=10>
<x=2, y=-7, z=3>
<x=9, y=-8, z=-3>`,
			numSteps:            100,
			expectedFinalEnergy: 1940,
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			i, e, g, c := day12.ProcessInput(tt.input)
			system := planetarySystem{i, e, g, c}

			t.Logf("After %d steps", 0)
			t.Logf("\n%v", system)

			for i := 0; i < tt.numSteps; i++ {
				system.step()
				t.Logf("After %d steps", i+1)
				t.Logf("\n%v", system)
			}

			totalEnergy := system.totalEnergy()
			if totalEnergy != tt.expectedFinalEnergy {
				t.Errorf("Expected energy of %d, got %d", tt.expectedFinalEnergy, totalEnergy)
			}
		})
	}
}

func TestStabalise(t *testing.T) {
	for i, tt := range []struct {
		input         string
		expectedSteps int
	}{
		{
			input: `<x=-1, y=0, z=2>
<x=2, y=-10, z=-7>
<x=4, y=-8, z=8>
<x=3, y=5, z=-1>`,
			expectedSteps: 2772,
		},
		{
			input: `<x=-8, y=-10, z=0>
<x=5, y=5, z=10>
<x=2, y=-7, z=3>
<x=9, y=-8, z=-3>`,
			expectedSteps: 4686774924,
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			i, e, g, c := day12.ProcessInput(tt.input)
			system := planetarySystem{i, e, g, c}

			steps := system.stepsUntilStable()

			if steps != tt.expectedSteps {
				t.Errorf("Expected steps of %d, got %d", tt.expectedSteps, steps)
			}
		})
	}
}
