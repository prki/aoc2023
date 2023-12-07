package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

var g_cardValueMap = map[string]int{
	"2": 2,
	"3": 3,
	"4": 4,
	"5": 5,
	"6": 6,
	"7": 7,
	"8": 8,
	"9": 9,
	"T": 10,
	"J": 11,
	"Q": 12,
	"K": 13,
	"A": 14,
}

var g_cardValueMapJoker = map[string]int{
	"2": 2,
	"3": 3,
	"4": 4,
	"5": 5,
	"6": 6,
	"7": 7,
	"8": 8,
	"9": 9,
	"T": 10,
	"J": 1,
	"Q": 12,
	"K": 13,
	"A": 14,
}

type Card struct {
	Label string
	Value int
}

type Hand struct {
	Cards []Card
	Bid   int
	Score int
}

type ByHandScore []Hand

func (b ByHandScore) Len() int      { return len(b) }
func (b ByHandScore) Swap(i, j int) { b[i], b[j] = b[j], b[i] }

func (b ByHandScore) Less(i, j int) bool {
	if b[i].Score != b[j].Score {
		return b[i].Score < b[j].Score
	}

	for cIdx := 0; cIdx < len(b[i].Cards); cIdx++ {
		if b[i].Cards[cIdx].Value < b[j].Cards[cIdx].Value {
			return true
		} else if b[i].Cards[cIdx].Value > b[j].Cards[cIdx].Value {
			return false
		}
	}

	return false
}

func (h *Hand) EvaluateScore() int {
	valueMap := make(map[int]int)
	for i := 0; i < len(h.Cards); i++ {
		valueMap[h.Cards[i].Value] += 1
	}

	threesCount := 0
	twosCount := 0
	for _, v := range valueMap {
		if v == 5 { // five-of-a-kind
			h.Score = 7
			return h.Score
		} else if v == 4 { // four-of-a-kind
			h.Score = 6
			return h.Score
		} else if v == 3 {
			threesCount += 1
		} else if v == 2 {
			twosCount += 1
		}
	}

	if threesCount == 1 && twosCount == 1 { // fullhouse
		h.Score = 5
	} else if threesCount == 1 && twosCount == 0 { // three-of-a-kind
		h.Score = 4
	} else if threesCount == 0 && twosCount == 2 { // two pair
		h.Score = 3
	} else if threesCount == 0 && twosCount == 1 { // pair
		h.Score = 2
	} else if threesCount == 0 && twosCount == 0 { // high card
		h.Score = 1
	} else {
		log.Fatal("[ERROR] Error evaluating score for hand", h)
	}

	return h.Score
}

// Brute force approach - tries all possible joker values, if any jokers
// are present.
// Proper approach should be designed around "joker" promotion - promote
// joker to the highest card in the hand based on the type of hand
// (e.g. 22AKJ should be 22AK2).
// For simplicity and ease of testing quickly, we do naive brute force.
func (h *Hand) EvaluateScoreWithJokers() int {
	isJokerInHand := false
	for i := 0; i < len(h.Cards); i++ {
		if h.Cards[i].Label == "J" {
			isJokerInHand = true
			break
		}
	}

	if !isJokerInHand {
		return h.EvaluateScore()
	}

	fmt.Println("Evaluating hand:", h)

	maxScore := 0
	for currJoker := 2; currJoker <= g_cardValueMapJoker["A"]; currJoker++ {
		tmpHand := Hand{}
		// manual copy because otherwise slice points to same array, overwriting h.Cards
		for i := 0; i < len(h.Cards); i++ {
			tmpHand.Cards = append(tmpHand.Cards, h.Cards[i])
		}

		for i := 0; i < len(tmpHand.Cards); i++ {
			if tmpHand.Cards[i].Label == "J" {
				tmpHand.Cards[i].Value = currJoker
			}
		}
		tmpHand.EvaluateScore()
		if tmpHand.Score > maxScore {
			maxScore = tmpHand.Score
			for i := 0; i < len(h.Cards); i++ {
				if h.Cards[i].Label == "J" {
					h.Cards[i].Value = currJoker
				}
			}
			h.Score = tmpHand.Score
		}
	}

	// [DBG]
	fmt.Println("Hand after promotion:", h)

	// Setting joker values to 1 so that comparison/sort behaves same as in solution1
	for i := 0; i < len(h.Cards); i++ {
		if h.Cards[i].Label == "J" {
			h.Cards[i].Value = 1
		}
	}

	return h.Score
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

func ParseHands(inputLines []string, isPartTwo bool) []Hand {
	var ret []Hand
	for i := 0; i < len(inputLines); i++ {
		hand := Hand{}
		lineFields := strings.Fields(inputLines[i])
		bid, err := strconv.Atoi(lineFields[1])
		if err != nil {
			log.Fatal("[ERROR] strconv err on bid:", err)
		}
		hand.Bid = bid
		for j := 0; j < 5; j++ {
			cardLabel := string(lineFields[0][j])
			value, ok := g_cardValueMap[cardLabel]
			if !ok {
				log.Fatal("[ERROR] g_cardValueMap did not have value for label", cardLabel)
			}
			card := Card{Label: cardLabel, Value: value}
			hand.Cards = append(hand.Cards, card)
		}

		if !isPartTwo {
			hand.EvaluateScore()
		} else {
			hand.EvaluateScoreWithJokers()
		}
		ret = append(ret, hand)
	}

	return ret
}

func Solution1(hands []Hand) uint64 {
	ret := uint64(0)

	sort.Sort(ByHandScore(hands))

	for i := 0; i < len(hands); i++ {
		rank := i + 1
		ret += uint64(rank) * uint64(hands[i].Bid)
		fmt.Println("Rank", i, hands[i])
	}

	return ret
}

func Solution2(hands []Hand) uint64 {
	return Solution1(hands)
}

func main() {
	inputLines := ReadInput("./input.txt")
	hands := ParseHands(inputLines, false)
	sol1 := Solution1(hands)
	fmt.Println("Solution 1:", sol1)

	hands2 := ParseHands(inputLines, true)
	sol2 := Solution2(hands2)
	fmt.Println("Solution 2:", sol2)
}
