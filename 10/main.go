package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// Coordinates in the tile map. Used as a unique ID for a pipe point.
type Coord struct {
	X uint
	Y uint
}

type PipeType byte

const (
	VERTICAL   PipeType = '|'
	HORIZONTAL          = '-'
	NORTHEAST           = 'L'
	NORTHWEST           = 'J'
	SOUTHWEST           = '7'
	SOUTHEAST           = 'F'
	START               = 'S'
)

// Naive quick implementation pretty printing a PipeType
func (pt *PipeType) PTReadable() string {
	switch *pt {
	case VERTICAL:
		return "VERTICAL"
	case HORIZONTAL:
		return "HORIZONTAL"
	case NORTHEAST:
		return "NORTHEAST"
	case NORTHWEST:
		return "NORTHWEST"
	case SOUTHWEST:
		return "SOUTHWEST"
	case SOUTHEAST:
		return "SOUTHEAST"
	case START:
		return "START"
	}
	fmt.Println("[WARN] Unable to pretty print pipetype which is string:", string(*pt), "byte:", (*pt))
	return "UNKNOWN"
}

type Node struct {
	Neighbors []*Node
	Type      PipeType
	Coords    Coord
}

func (n *Node) PrettyPrint() {
	fmt.Println("Node coords: X:", n.Coords.X, "Y:", n.Coords.Y, "type:", n.Type.PTReadable())
	fmt.Println("Neighbors on coords:")
	for i := 0; i < len(n.Neighbors); i++ {
		fmt.Println(*n.Neighbors[i])
	}
}

// To consider - the problem requires finding a strongly connected component
// of the graph where the animal is supposedly stuck. For now, I won't implement
// that, but who knows.
type Graph struct {
	Nodes   []*Node
	NodeMap map[Coord]*Node
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

func (n *Node) AddNeighbors(graph *Graph, crd1, crd2 Coord) {
	neighbor1, ok := graph.NodeMap[crd1]
	if !ok {
		//log.Fatal("Unable to get node on coords", crd1)
		log.Print("[WARN] Unable to get neighbor nodes for node pos:", n.Coords)
		return
	}
	neighbor2, ok := graph.NodeMap[crd2]
	if !ok {
		//log.Fatal("Unable to get node on coords", crd2)
		log.Print("[WARN] Unable to get neighbor nodes for node pos:", n.Coords)
		return
	}

	n.Neighbors = append(n.Neighbors, neighbor1, neighbor2)
	graph.NodeMap[n.Coords] = n
}

// Utility function calculating to-be neighbor coords based on a specific position
// and a pipe type.
func CalcNeighborCoords(pos Coord, pt PipeType) (Coord, Coord) {
	if pt == VERTICAL {
		neighborCoords1 := Coord{X: pos.X, Y: pos.Y - 1}
		neighborCoords2 := Coord{X: pos.X, Y: pos.Y + 1}
		return neighborCoords1, neighborCoords2
	} else if pt == HORIZONTAL {
		neighborCoords1 := Coord{X: pos.X + 1, Y: pos.Y}
		neighborCoords2 := Coord{X: pos.X - 1, Y: pos.Y}
		return neighborCoords1, neighborCoords2
	} else if pt == NORTHEAST {
		neighborCoords1 := Coord{X: pos.X, Y: pos.Y - 1} // NORTH
		neighborCoords2 := Coord{X: pos.X + 1, Y: pos.Y} // EAST
		return neighborCoords1, neighborCoords2
	} else if pt == NORTHWEST {
		neighborCoords1 := Coord{X: pos.X, Y: pos.Y - 1} // NORTH
		neighborCoords2 := Coord{X: pos.X - 1, Y: pos.Y} // WEST
		return neighborCoords1, neighborCoords2
	} else if pt == SOUTHEAST {
		neighborCoords1 := Coord{X: pos.X, Y: pos.Y + 1} // SOUTH
		neighborCoords2 := Coord{X: pos.X + 1, Y: pos.Y} // EAST
		return neighborCoords1, neighborCoords2
	} else if pt == SOUTHWEST {
		neighborCoords1 := Coord{X: pos.X, Y: pos.Y + 1} // SOUTH
		neighborCoords2 := Coord{X: pos.X - 1, Y: pos.Y} // WEST
		return neighborCoords1, neighborCoords2
	} else if pt == START {
		// nothing
	} else {
		log.Fatal("Unhandled node type when adding neighbors")
	}

	return Coord{}, Coord{} // should be unreachable
}

func (n *Node) HasNeighbor(other *Node) bool {
	for i := 0; i < len(n.Neighbors); i++ {
		//fmt.Println("neighbor[i] addr:", &n.Neighbors[i], "other addr:", &other)
		if n.Neighbors[i].Coords.X == other.Coords.X && n.Neighbors[i].Coords.Y == other.Coords.Y {
			return true
		}
	}
	return false
}

// Adds neighbors to start node. tries each option and checks whether 2 neighbors are present.
// If so, that particular option is the "true" node type.
func CalcStartNeighbors(graph *Graph, startNode *Node) {
	pipeTypes := []PipeType{VERTICAL, HORIZONTAL, NORTHEAST, NORTHWEST, SOUTHWEST, SOUTHEAST}
	for i := 0; i < len(pipeTypes); i++ {
		coord1, coord2 := CalcNeighborCoords(startNode.Coords, pipeTypes[i])
		node1, ok1 := graph.NodeMap[coord1]
		node2, ok2 := graph.NodeMap[coord2]
		if ok1 && ok2 {
			fmt.Println("Testing if start node is a neighbor of nodes:", node1, "and", node2)
			if node1.HasNeighbor(startNode) && node2.HasNeighbor(startNode) {
				startNode.Neighbors = append(startNode.Neighbors, node1, node2)
				return
			}
		}
	}
}

// Graph construction is done in a two-pass algorithm.
// In the first pass, we traverse the input to discover nodes and simply
// add them into the graph and nodemap. In the second pass, we traverse
// the constructed nodes and based on their type and coord, we add the
// possible edges (implemented as neighbor node pointer(s)).
func ConstructGraph(input []string) (*Graph, *Node) {
	ret := &Graph{
		Nodes:   []*Node{},
		NodeMap: make(map[Coord]*Node),
	}

	// First pass - create nodes
	for y := 0; y < len(input); y++ {
		for x := 0; x < len(input[y]); x++ {
			if input[y][x] == '.' {
				continue
			}
			coord := Coord{
				X: uint(x),
				Y: uint(y),
			}
			node := Node{
				//Neighbors: []*Node{},
				Neighbors: make([]*Node, 0, 2),
				Coords:    coord,
				Type:      PipeType(input[y][x]),
			}
			ret.Nodes = append(ret.Nodes, &node)
			ret.NodeMap[coord] = ret.Nodes[len(ret.Nodes)-1]
		}
	}

	// Second pass - assign neighbor nodes/edges
	var startNode *Node // neighbors defined only after all neighbors are set
	pipeTypes := []PipeType{VERTICAL, HORIZONTAL, NORTHEAST, NORTHWEST, SOUTHWEST, SOUTHEAST}
	for i := 0; i < len(ret.Nodes); i++ {
		currNode := ret.Nodes[i]
		if currNode.Type == START {
			startNode = currNode
		} else {
			for j := 0; j < len(pipeTypes); j++ {
				if currNode.Type == pipeTypes[j] {
					neigh1, neigh2 := CalcNeighborCoords(currNode.Coords, pipeTypes[j])
					//fmt.Println("Adding neighbors to coords:", currNode.Coords, neigh1, neigh2)
					currNode.AddNeighbors(ret, neigh1, neigh2)
				}
			}
		}
	}

	CalcStartNeighbors(ret, startNode)

	return ret, startNode
}

func BFSDistanceMap(graph *Graph, start *Node, distMap map[Coord]int) {
	/*fmt.Println("Building BFS Distance map for graph nodes:")
	for i := 0; i < len(graph.Nodes); i++ {
		fmt.Println(graph.Nodes[i])
	}
	*/
	visitedMap := make(map[Coord]bool)
	queue := []*Node{}
	queue = append(queue, start)
	visitedMap[start.Coords] = true
	distMap[start.Coords] = 0
	maxDist := 0
	for len(queue) > 0 {
		currNode := queue[0]
		queue = queue[1:] // dequeue
		//fmt.Println("BFS Dequeued node:", currNode)
		for i := 0; i < len(currNode.Neighbors); i++ {
			_, wasVisited := visitedMap[currNode.Neighbors[i].Coords]
			if !wasVisited {
				//fmt.Println("BFS CurrNode coords:", currNode.Coords, "Enqueued node:", currNode.Neighbors[i])
				queue = append(queue, currNode.Neighbors[i])
				//fmt.Println("Assigning neighbor", currNode.Neighbors[i].Coords, "distance", distMap[currNode.Coords], "+1")
				distMap[currNode.Neighbors[i].Coords] = distMap[currNode.Coords] + 1
				if distMap[currNode.Coords]+1 > maxDist {
					maxDist = distMap[currNode.Coords] + 1
				}
				visitedMap[currNode.Neighbors[i].Coords] = true
			}
		}
	}

	fmt.Println("Solution 1:", maxDist)
}

func Solution1(graph *Graph, start *Node) {
	distanceMap := make(map[Coord]int)
	BFSDistanceMap(graph, start, distanceMap)
	//fmt.Println("Distance map:", distanceMap)
}

func main() {
	input := ReadInput("input.txt")
	graph, startNode := ConstructGraph(input)

	Solution1(graph, startNode)
}
