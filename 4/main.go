package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type ScratchcardGame struct {
	WinningNumbers    map[int]int
	SelectedNumbers   []int
	PointsAwarded     int
	CntWinningNumbers int
	GameId            int
}

func (g *ScratchcardGame) CalculatePoints() {
	calcMap := make(map[int]int)
	for _, selNum := range g.SelectedNumbers {
		_, ok := g.WinningNumbers[selNum]
		if ok {
			calcMap[selNum] += 1
		}
	}

	winningCnt := 0
	for _, v := range calcMap {
		if v > 1 {
			fmt.Println("[WARN] Value was chosen more than once")
		}
		winningCnt += 1
	}

	if winningCnt == 0 {
		g.PointsAwarded = 0
	} else {
		g.PointsAwarded = int(math.Pow(2, float64(winningCnt-1))) // -1 - starting score is "1" == 2^0
	}

	g.CntWinningNumbers = winningCnt
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

func parseWinningNumbers(line string) map[int]int {
	ret := make(map[int]int)

	nums := strings.Split(line, " ")
	for _, num := range nums {
		winningNumber, err := strconv.Atoi(num)
		if err == nil {
			ret[winningNumber] += 1
		}
	}

	return ret
}

func parseSelectedNumbers(line string) []int {
	var ret []int

	nums := strings.Split(line, " ")
	for _, num := range nums {
		selectedNumber, err := strconv.Atoi(num)
		if err == nil {
			ret = append(ret, selectedNumber)
		}
	}

	return ret
}

func initializeGames(input []string) []ScratchcardGame {
	var ret []ScratchcardGame

	for _, line := range input {
		var game ScratchcardGame
		gameIdSplit := strings.Split(line, ":")
		gameIdStr := strings.Fields(gameIdSplit[0])[1]
		gameId, err := strconv.Atoi(gameIdStr)
		if err != nil {
			fmt.Println("[ERR] Error on atoi of gameidstr:", gameIdStr)
			log.Fatal(err)
		}
		game.GameId = gameId - 1 // gameid - 1 enables trivial array access as input array is ordered

		numberLine := strings.Split(gameIdSplit[1], "|")
		game.WinningNumbers = parseWinningNumbers(numberLine[0])
		game.SelectedNumbers = parseSelectedNumbers(numberLine[1])

		ret = append(ret, game)
	}

	return ret
}

func ParseScratchcardGames(input []string) []ScratchcardGame {
	games := initializeGames(input)
	for i := 0; i < len(games); i++ {
		games[i].CalculatePoints()
	}

	return games
}

func SpawnCopies(game ScratchcardGame, games []ScratchcardGame, copyMap map[int][]ScratchcardGame) {
	for i := 0; i < game.CntWinningNumbers; i++ {
		gameToCopy := i + game.GameId + 1
		if gameToCopy >= len(games) {
			return
		}

		entry, _ := copyMap[gameToCopy]
		entry = append(entry, games[gameToCopy])
		copyMap[gameToCopy] = entry
	}
}

// Suboptimal approach focused on copying - truly creates copies of games.
// Very inefficient both on time and memory.
// Likelier solution using the same approach would be to have a map[int]int
// where game id points to number of copies. Since we know the wincount, we
// could use that number to add +copycount to winning games.
func SolutionPartTwo(games []ScratchcardGame) int {
	cardCount := 0
	copyMap := make(map[int][]ScratchcardGame)

	for i := 0; i < len(games); i++ {
		SpawnCopies(games[i], games, copyMap)
		copiedGames, _ := copyMap[games[i].GameId]
		for j := 0; j < len(copiedGames); j++ {
			SpawnCopies(copiedGames[j], games, copyMap)
		}
	}

	cardCount += len(games)
	for k, _ := range copyMap {
		//fmt.Println("Len copymap", k, ":", len(copyMap[k]))
		cardCount += len(copyMap[k])
	}

	return cardCount
}

func main() {
	puzzleInput := ReadInput("./input.txt")
	games := ParseScratchcardGames(puzzleInput)
	//fmt.Println(games)
	solution := 0
	for i := 0; i < len(games); i++ {
		solution += games[i].PointsAwarded
	}

	fmt.Println("Solution 1:", solution)

	solution = SolutionPartTwo(games)
	fmt.Println("Solution 2:", solution)
}
