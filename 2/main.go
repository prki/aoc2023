package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type CubeGameRound struct {
	redCubes   int
	greenCubes int
	blueCubes  int
}

type CubeGame struct {
	id     int
	rounds []CubeGameRound
}

type CubeGameConfiguration struct {
	redCubesMax   int
	greenCubesMax int
	blueCubesMax  int
}

func ParseInputLine(inputLine string) CubeGame {
	var ret CubeGame
	gameIdSplit := strings.Split(inputLine, ":")
	var gameIdStr string = strings.Split(gameIdSplit[0], " ")[1]
	gameId, err := strconv.Atoi(gameIdStr)
	if err != nil {
		log.Fatal(err)
	}

	ret.id = gameId

	var gameRounds []string = strings.Split(gameIdSplit[1], ";")
	for _, rnd := range gameRounds {
		var gameRnd CubeGameRound
		var gameRndResult []string = strings.Split(rnd, ",")
		for _, res := range gameRndResult {
			cubePair := strings.Split(res, " ")
			cubeCnt, err := strconv.Atoi(cubePair[1]) // first char is always " " in input format
			if err != nil {
				log.Fatal(err)
			}
			cubeColor := cubePair[2]

			if cubeColor == "red" {
				gameRnd.redCubes += cubeCnt
			} else if cubeColor == "green" {
				gameRnd.greenCubes += cubeCnt
			} else if cubeColor == "blue" {
				gameRnd.blueCubes += cubeCnt
			} else {
				log.Fatal("Could not parse cube color", cubeColor)
			}
		}

		ret.rounds = append(ret.rounds, gameRnd)
	}

	return ret
}

func ValidateGame(game CubeGame, cfg CubeGameConfiguration) bool {
	for _, rnd := range game.rounds {
		if rnd.redCubes > cfg.redCubesMax {
			return false
		} else if rnd.greenCubes > cfg.greenCubesMax {
			return false
		} else if rnd.blueCubes > cfg.blueCubesMax {
			return false
		}
	}

	return true
}

func CalcMinimumConfiguration(game CubeGame) CubeGameConfiguration {
	var cfg CubeGameConfiguration

	for _, rnd := range game.rounds {
		if rnd.redCubes > cfg.redCubesMax {
			cfg.redCubesMax = rnd.redCubes
		}
		if rnd.greenCubes > cfg.greenCubesMax {
			cfg.greenCubesMax = rnd.greenCubes
		}
		if rnd.blueCubes > cfg.blueCubesMax {
			cfg.blueCubesMax = rnd.blueCubes
		}
	}

	return cfg
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

func SumValidIds(games []CubeGame) int {
	ret := 0

	for _, game := range games {
		ret += game.id
	}

	return ret
}

func main() {
	input := ReadInput("./input.txt")
	var cubeGames []CubeGame
	cfg := CubeGameConfiguration{
		redCubesMax:   12,
		greenCubesMax: 13,
		blueCubesMax:  14,
	}

	for _, lin := range input {
		cubeGames = append(cubeGames, ParseInputLine(lin))
	}

	var validCubeGames []CubeGame
	for _, game := range cubeGames {
		if ValidateGame(game, cfg) {
			validCubeGames = append(validCubeGames, game)
		}
	}

	// solution for first half
	/*solution := SumValidIds(validCubeGames)
	fmt.Println("Solution: ", solution)
	*/

	solution := 0
	for _, game := range cubeGames {
		smallestCfg := CalcMinimumConfiguration(game)
		power := smallestCfg.redCubesMax * smallestCfg.greenCubesMax * smallestCfg.blueCubesMax
		solution += power
	}

	fmt.Println("Solution 2nd half: ", solution)
}
