package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

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

// Replaces all string digits by their actual values (i.e. "one" -> "o1e").
// First and last characters are kept to enable shared characters, e.g. eightwo => becomes "e8tt2o".
// Since order of digits is first/last, even if something should not be considered a digit,
// we don't need to care.
func PreprocessLine(line string) string {
	transformationTable := [9][2]string{{"one", "o1e"}, {"two", "t2o"}, {"three", "t3e"}, {"four", "f4r"}, {"five", "f5e"}, {"six", "s6x"}, {"seven", "s7n"}, {"eight", "e8t"}, {"nine", "n9e"}}
	for i := 0; i < 9; i++ {
		line = strings.Replace(line, transformationTable[i][0], transformationTable[i][1], -1)
	}

	return line
}

// Naive straightforward implementation discovering digits char-by-char in-place
func ProcessLine(calibrationLine string) int {
	firstDigit := 0
	lastDigit := -1
	firstSet := false

	for i := 0; i < len(calibrationLine); i++ {
		digit, err := strconv.Atoi(string(calibrationLine[i]))
		if err == nil {
			if !firstSet {
				firstDigit = digit
				firstSet = true
			} else {
				lastDigit = digit
			}
		}
	}

	if lastDigit == -1 {
		lastDigit = firstDigit
	}

	ret, _ := strconv.Atoi(strconv.Itoa(firstDigit) + strconv.Itoa(lastDigit))

	return ret
}

func main() {
	fmt.Println("Hello, World!")

	input := ReadInput("./input.txt")
	sumResults := 0
	for _, line := range input {
		preprocessedLine := PreprocessLine(line)
		fmt.Println("Line processed: ", preprocessedLine)
		lineResult := ProcessLine(preprocessedLine)
		sumResults += lineResult
	}

	fmt.Println("Callibration result:", sumResults)
}
