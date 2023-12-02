package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	inputLines, err := extractInputLines("input.txt")
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

func extractInputLines(fileName string) ([]string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	inputfilePath := fmt.Sprintf("%s/%s", cwd, fileName)

	inputFile, err := os.Open(inputfilePath)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(inputFile)
	scanner.Split(bufio.ScanLines)

	var inputLines []string
	for scanner.Scan() {
		inputLines = append(inputLines, scanner.Text())
	}

	err = inputFile.Close()
	if err != nil {
		return nil, err
	}

	return inputLines, nil
}

func parseNumericCharsFromInput(input string, words *Trie) []int {
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

				if node := words.path(currentWord); node != nil {
					if *&node.wordEnds == true {
						numericChars = append(numericChars, *&node.val)
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

func getDigitWordTrie() *Trie {
	words := trie()
	words.insert("zero", 0)
	words.insert("one", 1)
	words.insert("two", 2)
	words.insert("three", 3)
	words.insert("four", 4)
	words.insert("five", 5)
	words.insert("six", 6)
	words.insert("seven", 7)
	words.insert("eight", 8)
	words.insert("nine", 9)

	return words
}

// **********************************************
// * START: Trie implementation
// **********************************************

type Node struct {
	children map[string]*Node
	val      int
	wordEnds bool
}

type Trie struct {
	root *Node
}

func trie() *Trie {
	t := new(Trie)
	t.root = new(Node)
	return t
}

func (t *Trie) insert(word string, val int) {
	current := t.root
	for _, cRune := range word {
		char := string(cRune)

		if current.children == nil {
			current.children = make(map[string]*Node)
		}

		if current.children[char] == nil {
			current.children[char] = new(Node)
		}
		current = current.children[char]
	}
	current.wordEnds = true
	current.val = val
}

func (t *Trie) path(partialWord string) *Node {
	current := t.root
	for i, cRune := range partialWord {
		childIdx := string(cRune)

		if i == 0 && current.children[childIdx] == nil {
			return nil
		}

		if current.children[childIdx] == nil {
			return current
		}

		current = current.children[childIdx]
	}

	return current
}
