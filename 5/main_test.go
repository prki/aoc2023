package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntervalOverlap(t *testing.T) {
	intA := Interval{
		Start: 2,
		End:   100,
	}
	intB := Interval{
		Start: 1,
		End:   6,
	}

	intActual := CalcIntervalOverlap(intA, intB)
	intExpected := Interval{
		Start: 2,
		End:   6,
	}

	assert.Equal(t, intExpected, intActual)
}

func TestIntervalNonoverlap(t *testing.T) {
	intA := Interval{
		Start: 2,
		End:   100,
	}
	intB := Interval{
		Start: 1,
		End:   6,
	}
	nonOverlaps := CalcIntervalNonoverlap(intA, intB)
	expected := []Interval{
		{
			Start: 1,
			End:   1,
		},
		{
			Start: 7,
			End:   100,
		},
	}

	assert.ElementsMatch(t, expected, nonOverlaps)

	nonOverlaps = CalcIntervalNonoverlap(intB, intA)
	assert.ElementsMatch(t, expected, nonOverlaps)
}
