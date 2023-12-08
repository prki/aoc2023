package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type TargetNode struct {
	LeftTarget  string
	RightTarget string
}

type FSM struct {
	Graph               map[string]TargetNode
	CurrStates          []string
	AcceptedInputsCount uint64
	InitStates          []string
	AcceptStates        []string
	AcceptCounter       map[string]int
}

// Initializes FSM by assigning the state graph and setting
// curr states to initial/input states.
func NewFSM(graph map[string]TargetNode) *FSM {
	var initStates []string
	var acceptStates []string
	for k := range graph {
		if k[len(k)-1] == 'A' {
			initStates = append(initStates, k)
		}
		if k[len(k)-1] == 'Z' {
			acceptStates = append(acceptStates, k)
		}
	}

	var currStatesCopy []string
	currStatesCopy = append(currStatesCopy, initStates...)

	fsm := &FSM{
		Graph:               graph,
		CurrStates:          currStatesCopy,
		InitStates:          initStates,
		AcceptStates:        acceptStates,
		AcceptedInputsCount: 0,
		AcceptCounter:       make(map[string]int),
	}

	return fsm
}

func (f *FSM) IsInAcceptState() bool {
	inAcceptState := true

	for i := 0; i < len(f.CurrStates); i++ {
		currState := f.CurrStates[i]
		if currState[len(currState)-1] == 'Z' {
			//fmt.Println("Init state:", f.InitStates[i], "is in accept state:", currState, "after", f.AcceptedInputsCount, "steps")
			f.AcceptCounter[currState] += 1
			if f.AcceptCounter[currState] == 1 {
				fmt.Println("Init state:", f.InitStates[i], "discovered accept state:", currState, "after", f.AcceptedInputsCount, "steps")
			} else if f.AcceptCounter[currState] == 2 {
				fmt.Println("Init state:", f.InitStates[i], "discovered same accept state:", currState, "after", f.AcceptedInputsCount, "steps")
			}
		}
		if currState[len(currState)-1] != 'Z' {
			inAcceptState = false
		}
	}

	return inAcceptState
}

// All of the current states accept input passed as a parameter
// and move to the next state.
func (f *FSM) AcceptInput(input byte) {
	for i := 0; i < len(f.CurrStates); i++ {
		if input == 'R' {
			f.CurrStates[i] = f.Graph[f.CurrStates[i]].RightTarget
		} else {
			f.CurrStates[i] = f.Graph[f.CurrStates[i]].LeftTarget
		}
	}
	f.AcceptedInputsCount += 1
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

func ParseGraph(input []string) (string, map[string]TargetNode) {
	graph := make(map[string]TargetNode)
	moveDirs := input[0]

	for i := 2; i < len(input); i++ {
		inputCut := strings.Replace(input[i], ",", "", -1)
		inputCut = strings.Replace(inputCut, "(", "", -1)
		inputCut = strings.Replace(inputCut, ")", "", -1)
		inputCut = strings.Replace(inputCut, "=", "", -1)
		nodes := strings.Fields(inputCut)
		startNode := nodes[0]
		target := TargetNode{
			LeftTarget:  nodes[1],
			RightTarget: nodes[2],
		}
		graph[startNode] = target
	}

	return moveDirs, graph
}

func FollowDirs(moveDirs string, graph map[string]TargetNode, startNode *string) int {
	if *startNode == "ZZZ" {
		return 0
	}

	stepCounter := 0
	currNode := *startNode
	for i := 0; i < len(moveDirs); i++ {
		stepCounter += 1
		if moveDirs[i] == 'R' {
			currNode = graph[currNode].RightTarget
		} else {
			currNode = graph[currNode].LeftTarget
		}

		if currNode == "ZZZ" {
			break
		}
	}

	*startNode = currNode

	return stepCounter
}

// "Naive"? approach which simply follows directions until target is found.
func Solution1(moveDirs string, graph map[string]TargetNode) int {
	totalSteps := 0
	lastNode := "AAA"
	for {
		steps := FollowDirs(moveDirs, graph, &lastNode)
		totalSteps += steps
		if lastNode == "ZZZ" {
			break
		}
	}

	return totalSteps
}

// Through analyzing the graph, it became clear the graph is actually
// consisting of count(init_states) subgraphs. These graphs are then
// looping in a period.
// Therefore, assuming that e.g. one subgraph has a period of 2 and another
// of 5, we are in accepting states after 10 steps. This is equivalent to
// finding LCM.
// That being said, in this case, surprisingly enough the period itself
// turned out to be the same as distance from the start.
// For ease of development, this was executed to provide all the distances
// + periods and then were just put into an online LCM calculator.
func Solution2(moveDirs string, graph map[string]TargetNode) uint64 {
	totalSteps := uint64(0)
	fsm := NewFSM(graph)
	moveDirIdx := 0
	for {
		totalSteps += 1
		if moveDirIdx >= len(moveDirs) {
			moveDirIdx = 0
		}
		fsm.AcceptInput(moveDirs[moveDirIdx])
		isAccepted := fsm.IsInAcceptState()
		foundCounter := 0
		for _, v := range fsm.AcceptCounter {
			if v >= 2 {
				foundCounter += 1
			}
		}
		if foundCounter == len(fsm.InitStates) {
			fmt.Println("All init states with periods of movement found.")
			return totalSteps
		}
		if isAccepted {
			break
		}
		moveDirIdx++
	}

	return totalSteps
}

func main() {
	inputLines := ReadInput("./input.txt")
	moveDirs, graph := ParseGraph(inputLines)
	/*sol1 := Solution1(moveDirs, graph)
	fmt.Println("Solution 1:", sol1)
	*/

	sol2 := Solution2(moveDirs, graph)
	fmt.Println("Solution 2:", sol2)
}
