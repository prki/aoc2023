package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type BoatRace struct {
	RecordTime     int
	RecordDistance int
}

type Simulation struct {
	TimeHeld        int
	Velocity        int
	DistanceCovered int
}

func ReadInput(path string) []string {
	ret := make([]string, 0)
	fil, err := os.Open(path)
	if err != nil {
		log.Fatalf("[ERROR] Can't open file %s\n", path)
	}
	defer fil.Close()

	scanner := bufio.NewScanner(fil)
	for scanner.Scan() {
		ret = append(ret, scanner.Text())
	}

	return ret
}

func ParseInput(inputLines []string) []BoatRace {
	var ret []BoatRace
	times := strings.Fields(inputLines[0])
	distances := strings.Fields(inputLines[1])

	for ti := 1; ti < len(times); ti++ {
		time, err := strconv.Atoi(times[ti])
		if err != nil {
			log.Fatal("[ERROR] Err strconv time parsing", err)
		}
		distance, err := strconv.Atoi(distances[ti])
		if err != nil {
			log.Fatal("[ERROR] Err strconv distance parsing", err)
		}
		ret = append(ret, BoatRace{RecordTime: time, RecordDistance: distance})
	}

	return ret
}

func ParseInput2(inputLines []string) BoatRace {
	var ret BoatRace
	timeFields := strings.Fields(inputLines[0])
	timeSubStrs := strings.Join(timeFields[1:], "")
	time, err := strconv.Atoi(timeSubStrs)
	if err != nil {
		log.Fatal("[ERROR] Err strconv time parsing", err)
	}

	distFields := strings.Fields(inputLines[1])
	distSubStrs := strings.Join(distFields[1:], "")
	dist, err := strconv.Atoi(distSubStrs)
	if err != nil {
		log.Fatal("[ERROR] Err strconv dist parsing", err)
	}

	ret.RecordTime = time
	ret.RecordDistance = dist

	return ret
}

func SimulateBoatrace(boatRace BoatRace) []Simulation {
	var ret []Simulation

	for i := 0; i < boatRace.RecordTime; i++ {
		timeHeld := i
		velocity := timeHeld * 1
		timeLeft := boatRace.RecordTime - timeHeld
		distance := velocity * timeLeft
		ret = append(ret, Simulation{TimeHeld: timeHeld, Velocity: velocity, DistanceCovered: distance})
	}

	return ret
}

// Naive solution which creates simulation for each time + hold option.
func Solution1(boatRaces []BoatRace) int {
	product := 1
	for _, boatRace := range boatRaces {
		simulations := SimulateBoatrace(boatRace)
		waysToWin := 0
		for _, simul := range simulations {
			if simul.DistanceCovered > boatRace.RecordDistance {
				waysToWin += 1
			}
		}
		if waysToWin > 1 {
			product *= waysToWin
		}
	}

	return product
}

// Naive solution creating simulations.
// Proper solution should consider time as a function and analyze where the function is >= dist.
// T_race - const
// T_held
// T_move
// T_race = T_held + T_move
// T_move = T_race - T_held
// v = T_held
// => s = v*t
// => s = t_held * t_move
// => s = (T_race - T_move) * T_move
// since s (distance) has a set minimum, we can consider this as a problem of solving
// the quadratic inequivality:
// minDist < (T_race - T_move) * T_move
// which would give us the the interval where the parabola has greater values than
// the distance.
// However, there is no reason to program this really, as the runtime is sufficient even
// with a "naive" solution :)
func Solution2(boatRace BoatRace) int {
	fmt.Println("Solution 2 boatrace:", boatRace)
	waysToWin := 0

	simulations := SimulateBoatrace(boatRace)
	for i := 0; i < len(simulations); i++ {
		if simulations[i].DistanceCovered > boatRace.RecordDistance {
			waysToWin += 1
		}
	}

	return waysToWin
}

func main() {
	input := ReadInput("./input.txt")
	boatRaces := ParseInput(input)
	sol1 := Solution1(boatRaces)
	fmt.Println("Solution 1:", sol1)

	boatRace2 := ParseInput2(input)
	sol2 := Solution2(boatRace2)
	fmt.Println("Solution 2:", sol2)
}
