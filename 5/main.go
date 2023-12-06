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

type Interval struct {
	Start uint64
	End   uint64
}

type MappingLine struct {
	SourceStart      uint64
	DestinationStart uint64
	Range            uint64
}

// To consider - putting a pointer to the next almanac map?
type AlmanacMap struct {
	Map      map[uint64]uint64 // inefficient due to massive ranges - but viable for testing purposes
	Mappings []MappingLine
}

func (am *AlmanacMap) PopulateMap() {
	for i := 0; i < len(am.Mappings); i++ {
		for j := uint64(0); j < am.Mappings[i].Range; j++ {
			srcStart := am.Mappings[i].SourceStart
			dstStart := am.Mappings[i].DestinationStart
			am.Map[srcStart+j] = dstStart + j
		}
	}
}

func (am *AlmanacMap) GetDestinationID_naive(source uint64) uint64 {
	ret, ok := am.Map[source]
	if ok {
		return ret
	} else {
		return source
	}
}

func (am *AlmanacMap) GetDestinationID(source uint64) uint64 {
	ret := source
	for i := 0; i < len(am.Mappings); i++ {
		srcStart := am.Mappings[i].SourceStart
		dstStart := am.Mappings[i].DestinationStart
		mapRange := am.Mappings[i].Range

		// number in range
		if source >= srcStart && source < srcStart+mapRange {
			tmp := source - srcStart
			ret = dstStart + tmp
			return ret
		}
	}

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

func parseSeeds(line string) []uint64 {
	var ret []uint64

	seedIds := strings.Fields(strings.Split(line, ":")[1])
	for i := 0; i < len(seedIds); i++ {
		//idInt, err := strconv.Atoi(seedIds[i])
		idInt, err := strconv.ParseUint(seedIds[i], 10, 64)
		if err != nil {
			log.Fatal("[ERROR] Cant parse seeds:", err)
		}
		ret = append(ret, idInt)
	}

	return ret
}

func parseMappingLine(inputLine string) MappingLine {
	fields := strings.Fields(inputLine)
	//destId, err := strconv.Atoi(fields[0])
	destId, err := strconv.ParseUint(fields[0], 10, 64)
	if err != nil {
		log.Fatal("error parsing mapping line:", err)
	}
	//srcId, err := strconv.Atoi(fields[1])
	srcId, err := strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		log.Fatal("error parsing mapping line:", err)
	}
	//mapRange, err := strconv.Atoi(fields[2])
	mapRange, err := strconv.ParseUint(fields[2], 10, 64)
	if err != nil {
		log.Fatal("error parsing mapping line:", err)
	}

	/*if destId < 0 || srcId < 0 || mapRange < 0 {
		fmt.Println("ERROR One of numbers could not fit into int")
		fmt.Println("Dest id:", destId, "src ID:", srcId, "Map range:", mapRange)
	}
	*/

	ret := MappingLine{
		SourceStart:      srcId,
		DestinationStart: destId,
		Range:            mapRange,
	}

	return ret
}

func ParseInput(inputLines []string) ([]uint64, []AlmanacMap) {
	almanacMaps := make([]AlmanacMap, 7, 7)
	for i := 0; i < len(almanacMaps); i++ {
		almanacMaps[i].Map = make(map[uint64]uint64)
	}
	var seeds []uint64
	almanacIdx := 0

	seeds = parseSeeds(inputLines[0])
	for i := 1; i < len(inputLines); i++ {
		if strings.Contains(inputLines[i], "map:") {
			for {
				i += 1
				if i >= len(inputLines) || len(inputLines[i]) == 0 { // newline - len == 0 it seems
					almanacIdx += 1
					break
				}
				mappingLine := parseMappingLine(inputLines[i])
				almanacMaps[almanacIdx].Mappings = append(almanacMaps[almanacIdx].Mappings, mappingLine)
			}
		}
	}

	// Uncomment for naive solution enablement
	/*for i := 0; i < len(almanacMaps); i++ {
		fmt.Println("populating map")
		almanacMaps[i].PopulateMap()
	}
	*/

	return seeds, almanacMaps
}

func MaxU64(a, b uint64) uint64 {
	if a >= b {
		return a
	}
	return b
}

func MinU64(a, b uint64) uint64 {
	if a <= b {
		return a
	}

	return b
}

func CalcIntervalOverlap(intA, intB Interval) Interval {
	ret := Interval{}

	if intA.End < intB.Start || intB.End < intA.Start {
		return ret
	}

	ret.Start = MaxU64(intA.Start, intB.Start)
	ret.End = MinU64(intA.End, intB.End)

	return ret
}

func calcIntervalNonoverlap(intA, intB Interval) []Interval {
	var ret []Interval
	// no overlap - just return both intervals
	if intA.End < intB.Start || intB.End < intA.Start {
		return []Interval{intA, intB}
	}

	if intA.Start < intB.Start {
		ret = append(ret, Interval{intA.Start, intB.Start - 1})
	}
	if intA.End > intB.End {
		ret = append(ret, Interval{intB.End + 1, intA.End})
	}

	return ret
}

func CalcIntervalNonoverlap(intA, intB Interval) []Interval {
	nonOverlapAB := calcIntervalNonoverlap(intA, intB)
	nonOverlapBA := calcIntervalNonoverlap(intB, intA)
	ret := append(nonOverlapAB, nonOverlapBA...)

	return ret
}

func MapToDestinationRange(srcInt Interval, mapping MappingLine) Interval {
	var diff int64 = int64(mapping.DestinationStart) - int64(mapping.SourceStart) // transformation rule diff
	var start int64 = int64(srcInt.Start) + diff
	var srcIntLen = int64(srcInt.End) - int64(srcInt.Start) + 1
	var end int64 = start + srcIntLen - 1
	ret := Interval{
		Start: uint64(start),
		End:   uint64(end),
	}

	//fmt.Println("Source interval", srcInt, "mapped to:", ret)

	return ret
}

func FilterOverlaps(int Interval, ints []Interval) []Interval {
	//fmt.Println("Filtering interval", int, "from overlap subsets", ints)
	toFilter := []Interval{int}
	for i := 0; i < len(ints); i++ {
		var tmp []Interval
		for j := 0; j < len(toFilter); j++ {
			nonOverlaps := CalcIntervalNonoverlap(toFilter[j], ints[i])
			tmp = append(tmp, nonOverlaps...)
		}
		toFilter = nil                      // clean old state to only work with not-yet-overlapping
		toFilter = append(toFilter, tmp...) // copy new nonoverlaps to tofilter
		tmp = nil
	}

	//fmt.Println("Filtered:", toFilter)

	return toFilter
}

func GeneratePossibleDestIntervals(srcInterval Interval, almanacMap AlmanacMap) []Interval {
	var ret []Interval
	var discoveredOverlaps []Interval

	for i := 0; i < len(almanacMap.Mappings); i++ {
		almanacSourceRange := Interval{
			Start: almanacMap.Mappings[i].SourceStart,
			End:   almanacMap.Mappings[i].SourceStart + almanacMap.Mappings[i].Range - 1,
		}
		overlap := CalcIntervalOverlap(srcInterval, almanacSourceRange)
		if overlap.Start == 0 && overlap.End == 0 {
			// source range just points to same destination values
		} else {
			discoveredOverlaps = append(discoveredOverlaps, overlap)
			destRange := MapToDestinationRange(overlap, almanacMap.Mappings[i])
			ret = append(ret, destRange) // nonoverlap just points to the same points
		}
	}

	directMappings := FilterOverlaps(srcInterval, discoveredOverlaps)
	ret = append(ret, directMappings...)

	return ret
}

// Solution for part 2 is built around interval arithmetics.
// The idea is that each source interval is mapped to a destination interval.
// In case there are overlap between the transformation domain (possible source inputs)
// and source intervals, new target intervals are calculated. If not, same intervals
// are reused.
// As such, all possible destination intervals are generated with respect to a certain
// interval input.
func Solution2(seeds []uint64, almanacMaps []AlmanacMap) uint64 {
	ret := uint64(0)
	var seedIntervals []Interval
	for i := 0; i < len(seeds); i += 2 {
		tmp := Interval{
			Start: seeds[i],
			End:   seeds[i] + seeds[i+1] - 1, // -1 - 2, 3 -> [2 3 4]
		}

		seedIntervals = append(seedIntervals, tmp)
	}

	var destIntervals []Interval
	currSrcInts := seedIntervals
	for i := 0; i < len(almanacMaps); i++ {
		for j := 0; j < len(currSrcInts); j++ {
			//fmt.Println("Generating possible intervals for source int:", currSrcInts[j])
			tmp := GeneratePossibleDestIntervals(currSrcInts[j], almanacMaps[i])
			destIntervals = append(destIntervals, tmp...)
		}
		currSrcInts = nil
		currSrcInts = append(currSrcInts, destIntervals...)
		destIntervals = nil
		//fmt.Println("Discovered dest intervals:", currSrcInts)
	}

	//fmt.Println("Dest intervals:", currSrcInts)
	ret = math.MaxUint64
	for i := 0; i < len(currSrcInts); i++ {
		if currSrcInts[i].Start < ret {
			ret = currSrcInts[i].Start
		}
	}

	return ret
}

func main() {
	input := ReadInput("./input.txt")
	seeds, almanacMaps := ParseInput(input)
	fmt.Println("Input parsed successfully")
	fmt.Println("Seeds:", seeds)
	var minLocation uint64 = math.MaxUint64
	for _, seed := range seeds {
		//fmt.Println("Seed mapping for seed:", seed)
		//fmt.Println("----------")
		//fmt.Println("Calculating seed mappings for seed:", seed)
		tmp := seed
		for i := 0; i < len(almanacMaps); i++ {
			tmp = almanacMaps[i].GetDestinationID(tmp)
			//fmt.Println("Almanac map pointed to id:", tmp)
		}
		if tmp < minLocation {
			minLocation = tmp
		}
	}

	fmt.Println("Solution 1:", minLocation)

	solution2 := Solution2(seeds, almanacMaps)
	fmt.Println("Solution 2:", solution2)
}
