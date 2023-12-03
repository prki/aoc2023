package main

/*Solution to 3 is built by creating a bounding box around all numbers to check
symbols in. If there is any symbol in this bounding box other than '.', we know
it's a "part number".

For part 2, a hashmap coord -> (int, int) is created, where first int in the pair
contains number of adjacent numbers to an asterisk at a coord, second value
is the product of number values.

Other interesting solutions considered included creating a graph where nodes
would be numbers/symbols and an edge would exist between them if they
are adjacent.
*/

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

// Utility structure for working with coords in schematic.
// schematic[y][x]
type Coord struct {
	Y int
	X int
}

type IntPair struct {
	First  int
	Second int
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

// dx = length of any string
// dy = string count
func LoadSchematic(diagram []string) [][]byte {
	dx := len(diagram[0])
	dy := len(diagram)

	ret := make([][]byte, dy)
	for i := range ret {
		ret[i] = make([]byte, dx)
	}

	for i := 0; i < len(diagram); i++ {
		for j := 0; j < len(diagram[i]); j++ {
			ret[i][j] = diagram[i][j]
		}
	}

	return ret
}

// Create a bounding box, check if bounding box contains non-digit/non-dot
// E.g.
//
//	1 2 3 4 5
//
// 1. . . . .
// 2. | - - |
// 3. | 6 9 |
// 4. | - - |
// 5. . . . .
// Number takes [3][3] -> [3][4]. Number is only in a row (3-4)
// Bounding points [2][2] [2][3] [2][4] [2][5]
// Bounding points [3][2]               [3][5]
// Bounding points [4][2] [4][3] [4][4] [4][5]
// May return illegal coordinates - meaning coordinates which are out of bounds
func CreateBoundingBox(numStart Coord, numEnd Coord) []Coord {
	var ret []Coord

	for i := numStart.X - 1; i <= numEnd.X; i++ {
		for j := numStart.Y - 1; j <= numStart.Y+1; j++ {
			var boxCoord Coord = Coord{
				X: i,
				Y: j,
			}
			ret = append(ret, boxCoord)
		}
	}

	return ret
}

// In case of invalid coords (out of bounds), return 0,
// otherwise return value
func GetSchematicValue(schematic [][]byte, coord Coord) byte {
	if coord.Y >= len(schematic) || coord.Y < 0 {
		return 0
	} else if coord.X >= len(schematic[coord.Y]) || coord.X < 0 {
		return 0
	}

	return schematic[coord.Y][coord.X]
}

func IsSymbolInBoundingBox(schematic [][]byte, boxCoords []Coord) (bool, string, Coord) {
	var retStr string = ""
	var retCoord Coord

	for _, c := range boxCoords {
		val := GetSchematicValue(schematic, c)
		if val != 0 {
			valStr := string(val)
			_, notNumErr := strconv.Atoi(valStr)
			if valStr != "." && notNumErr != nil {
				retStr = valStr
				retCoord = c
				return true, retStr, retCoord
			}
		}
	}

	return false, "", retCoord
}

func PrintSchematic(schematic [][]byte) {
	for y := range schematic {
		fmt.Println(string(schematic[y]))
	}
}

func FindNumberInSchematic(schematic [][]byte, numStart Coord) (int, Coord) {
	var retNum int
	var retCoord Coord
	var numBytes []byte

	retCoord.Y = numStart.Y

	for x := numStart.X; x < len(schematic[numStart.Y]); x++ {
		_, err := strconv.Atoi(string(schematic[numStart.Y][x]))
		if err != nil {
			retCoord.X = x
			break
		}
		numBytes = append(numBytes, schematic[numStart.Y][x])
		retCoord.X = x
	}

	retNum, _ = strconv.Atoi(string(numBytes))

	return retNum, retCoord
}

func RegisterAsterisk(asteriskMap map[Coord]IntPair, coord Coord, partNumber int) {
	val, ok := asteriskMap[coord]
	if !ok {
		val.Second = 1 // so that we dont multiply by 0
	}
	val.First += 1
	val.Second *= partNumber
	asteriskMap[coord] = val
}

func main() {
	diagram := ReadInput("input.txt")
	schematic := LoadSchematic(diagram)
	PrintSchematic(schematic)
	asteriskMap := make(map[Coord]IntPair)

	solutionSum := 0
	for y := 0; y < len(schematic); y++ {
		for x := 0; x < len(schematic[y]); x++ {
			_, err := strconv.Atoi(string(schematic[y][x]))
			if err == nil {
				numStartCoord := Coord{
					X: x,
					Y: y,
				}
				num, numEndCoord := FindNumberInSchematic(schematic, numStartCoord)
				boundingBoxCoords := CreateBoundingBox(numStartCoord, numEndCoord)
				symbolFound, symbol, coord := IsSymbolInBoundingBox(schematic, boundingBoxCoords)
				if symbolFound {
					solutionSum += num
				}
				if symbol == "*" {
					RegisterAsterisk(asteriskMap, coord, num)
				}
				x = numEndCoord.X
			}
		}
	}

	solution2 := 0
	for _, v := range asteriskMap {
		if v.First == 2 {
			solution2 += v.Second
		}
	}

	fmt.Println("Solution 1:", solutionSum)
	fmt.Println("Solution 2:", solution2)
}
