package day1_2023

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/Drew-Kimberly/advent-of-code/src/common/go/fs"
	"github.com/Drew-Kimberly/advent-of-code/src/common/go/trie"
)

func Day1_2023() {
	inputPath, err := filepath.Abs("./2023/day1/input.txt")
	if err != nil {
		panic(err)
	}

	inputLines, err := fs.ExtractInputLines(inputPath)
	if err != nil {
		panic(err)
	}

	totalSumPt1 := 0
	totalSumPt2 := 0

	words := getDigitWordTrie()

	for _, line := range inputLines {
		totalSumPt1 += calculateCalibrationValue(parseNumericCharsFromInput(line, nil))
		totalSumPt2 += calculateCalibrationValue(parseNumericCharsFromInput(line, words))
	}

	fmt.Println(fmt.Sprintf("Part 1 value: %d", totalSumPt1))
	fmt.Println(fmt.Sprintf("Part 2 value: %d", totalSumPt2))
}

func parseNumericCharsFromInput(input string, words *trie.ITrie) []int {
	var numericChars []int
	i := 0

	for i < len(input) {
		char := string(input[i])
		num, err := strconv.Atoi(char)
		if err == nil {
			numericChars = append(numericChars, num)
			i++
			continue
		}

		if words != nil {
			currentWord := char
			j := i

			for j < len(input) {
				if j > i {
					currentWord += string(input[j])
				}

				if node := words.Path(currentWord); node != nil {
					if *&node.WordEnds == true {
						numericChars = append(numericChars, *&node.Val)
						j = len(input)
						// Note - CANNOT set i=j since we'd miss cases like "4twone"
					}
					j++
				} else {
					// Sequence of chars cannot possibly a numeric word.
					j = len(input)
				}
			}
		}

		i++
	}

	return numericChars
}

func calculateCalibrationValue(lineDigits []int) int {
	var calibrationValue int
	var err error

	if len(lineDigits) < 1 {
		panic("At least 1 numeric characters expected per input line")
	} else if len(lineDigits) == 1 {
		calibrationValue, err = strconv.Atoi(fmt.Sprintf("%d%d", lineDigits[0], lineDigits[0]))
		if err != nil {
			panic(err)
		}
	} else {
		calibrationValue, err = strconv.Atoi(fmt.Sprintf("%d%d", lineDigits[0], lineDigits[len(lineDigits)-1]))
		if err != nil {
			panic(err)
		}
	}

	return calibrationValue
}

func getDigitWordTrie() *trie.ITrie {
	words := trie.Trie()
	words.Insert("zero", 0)
	words.Insert("one", 1)
	words.Insert("two", 2)
	words.Insert("three", 3)
	words.Insert("four", 4)
	words.Insert("five", 5)
	words.Insert("six", 6)
	words.Insert("seven", 7)
	words.Insert("eight", 8)
	words.Insert("nine", 9)

	return words
}
