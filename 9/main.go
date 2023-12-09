package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Struct "abstracting" the calculation tree.
// Contains slices/arrays of all the diff steps together with a New method
// which executes the calculation.
// Not implemented as a tree, though part 1 hints at a tree being valuable, but
// let's see.
type HistoryDiffs struct {
	Histories [][]int64
}

func NewHistoryDiffs(firstHist []int64) *HistoryDiffs {
	hists := make([][]int64, 1)
	hists[0] = append(hists[0], firstHist...)

	ret := &HistoryDiffs{
		Histories: hists,
	}

	return ret
}

// Stops when a slice [0, 0, ..., 0] is pushed into diffs
func (h *HistoryDiffs) CalculateHistoryDiffs() {
	stopCondition := true
	for currHistoryIdx := 0; ; currHistoryIdx++ {
		var newHistory []int64
		for j := 0; j < len(h.Histories[currHistoryIdx])-1; j++ { // -1 - going up to before-last elem
			left := h.Histories[currHistoryIdx][j]
			right := h.Histories[currHistoryIdx][j+1]
			diff := right - left
			if diff != 0 {
				stopCondition = false
			}
			newHistory = append(newHistory, diff)
		}
		h.Histories = append(h.Histories, newHistory)
		if stopCondition {
			return
		}
		stopCondition = true
	}
}

// Function extrapolating the next value for a history.
// It is expected that CalculateHistoryDiffs() has been called and
// the history is fully populated based on the algorithm defined in
// problem definition. Other calls are unsafe.
// This method also modifies the underlying histories.
func (h *HistoryDiffs) ExtrapolateValue() int64 {
	lastRow := len(h.Histories) - 1
	h.Histories[lastRow] = append(h.Histories[lastRow], 0)

	for row := len(h.Histories) - 1; row > 0; row-- {
		child := h.Histories[row][len(h.Histories[row])-1]
		leftParent := h.Histories[row-1][len(h.Histories[row-1])-1]
		rightParent := leftParent + child
		h.Histories[row-1] = append(h.Histories[row-1], rightParent)
	}

	ret := h.Histories[0][len(h.Histories[0])-1]
	return ret
}

// Function for solving part 2. Refer to docstring of HistoryDiffs.ExtrapolateValue()
// for more details, as the implementation strategy is equivalent.
func (h *HistoryDiffs) ExtrapolateValuePast() int64 {
	lastRow := len(h.Histories) - 1
	h.Histories[lastRow] = append([]int64{0}, h.Histories[lastRow]...) // prepend

	for row := len(h.Histories) - 1; row > 0; row-- {
		child := h.Histories[row][0]
		right := h.Histories[row-1][0]
		left := right - child
		h.Histories[row-1] = append([]int64{left}, h.Histories[row-1]...) // prepend
	}

	ret := h.Histories[0][0]
	return ret
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

func ParseHistories(input []string) [][]int64 {
	ret := make([][]int64, len(input))
	for i := 0; i < len(input); i++ {
		line := input[i]
		nums := strings.Fields(line)
		for j := 0; j < len(nums); j++ {
			currNum, err := strconv.ParseInt(nums[j], 10, 64)
			if err != nil {
				log.Fatal("[ERROR] err parsing histories:", err)
			}
			ret[i] = append(ret[i], currNum)
		}
	}

	return ret
}

func Solution1(hists [][]int64) int64 {
	ret := int64(0)

	for i := 0; i < len(hists); i++ {
		histDiffs := NewHistoryDiffs(hists[i])
		histDiffs.CalculateHistoryDiffs()
		/*fmt.Println("[DBG] History calculation for", hists[i])
		for j := 0; j < len(histDiffs.Histories); j++ {
			fmt.Println("Hist step", j, histDiffs.Histories[j])
		}
		*/

		extrapolated := histDiffs.ExtrapolateValue()
		//fmt.Println("Extrapolated value for hist", i, "-", extrapolated)
		ret += extrapolated
	}

	return ret
}

func Solution2(hists [][]int64) int64 {
	ret := int64(0)

	for i := 0; i < len(hists); i++ {
		histDiffs := NewHistoryDiffs(hists[i])
		histDiffs.CalculateHistoryDiffs()

		extrapolated := histDiffs.ExtrapolateValuePast()
		ret += extrapolated
	}

	return ret
}

func main() {
	input := ReadInput("input.txt")
	histories := ParseHistories(input)
	sol1 := Solution1(histories)
	fmt.Println("Solution 1:", sol1)

	sol2 := Solution2(histories)
	fmt.Println("Solution 2:", sol2)
}
