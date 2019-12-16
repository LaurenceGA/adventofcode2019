package main

import "testing"

import "strconv"
import day14 "github.com/LaurenceGA/adventofcode2019/14"

func TestOreToFuel(t *testing.T) {
	for i, tt := range []struct {
		input       string
		expectedOre int
	}{
		{
			input: `10 ORE => 10 A
1 ORE => 1 B
7 A, 1 B => 1 C
7 A, 1 C => 1 D
7 A, 1 D => 1 E
7 A, 1 E => 1 FUEL`,
			expectedOre: 31,
		},
		{
			input: `9 ORE => 2 A
8 ORE => 3 B
7 ORE => 5 C
3 A, 4 B => 1 AB
5 B, 7 C => 1 BC
4 C, 1 A => 1 CA
2 AB, 3 BC, 4 CA => 1 FUEL`,
			expectedOre: 165,
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			formulae := day14.ProcessInput(tt.input)
			requiredOre := amountOfOreToProduceFuel(formulae)
			if requiredOre != tt.expectedOre {
				t.Errorf("Expected %d, got %d", tt.expectedOre, requiredOre)
			}
		})
	}
}
