package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLine(t *testing.T) {
	tin := "Game 5: 2 red, 3 blue; 4 green, 10 red, 20 blue"
	var expectedRnds []CubeGameRound = []CubeGameRound{
		{
			redCubes:   2,
			blueCubes:  3,
			greenCubes: 0,
		},
		{
			redCubes:   10,
			greenCubes: 4,
			blueCubes:  20,
		},
	}
	expectedGame := CubeGame{
		id:     5,
		rounds: expectedRnds,
	}
	actual := ParseInputLine(tin)
	assert.Equal(t, expectedGame, actual)
}
