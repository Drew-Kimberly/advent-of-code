package day4_2023

import (
	"fmt"
	"math"
	"path/filepath"
	"strings"

	"github.com/Drew-Kimberly/advent-of-code/src/common/go/convert"
	"github.com/Drew-Kimberly/advent-of-code/src/common/go/fs"
	"github.com/Drew-Kimberly/advent-of-code/src/common/go/list"
	"github.com/yourbasic/bit"
)

const DELIMITER = " | "
const LINE_PREFIX = ":"

type ScratchCard struct {
	Numbers        []int
	WinningNumbers *bit.Set
}

func Day4_2023() {
	inputPath, err := filepath.Abs("./2023/day4/input.txt")
	if err != nil {
		panic(err)
	}

	inputLines, err := fs.ExtractInputLines(inputPath)
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("Part 1 value: %d", calculatePoints(parseScratchCards(inputLines))))
	fmt.Println(fmt.Sprintf("Part 2 value: %d", calculateTotalScratchCards(parseScratchCards(inputLines))))
}

func calculatePoints(scratchCards []*ScratchCard) int {
	points := 0
	for _, card := range scratchCards {
		numWinning := card.NumberOfWins()
		if numWinning > 0 {
			points += int(math.Pow(2, float64(numWinning-1)))
		}
	}

	return points
}

func calculateTotalScratchCards(originalScratchCards []*ScratchCard) int {
	total := len(originalScratchCards)

	for i, card := range originalScratchCards {
		numWinning := card.NumberOfWins()
		total += calculateTotalScratchCards(originalScratchCards[i+1 : i+numWinning+1])
	}

	return total
}

func parseScratchCards(inputLines []string) []*ScratchCard {
	var scratchCards []*ScratchCard

	for _, line := range inputLines {
		prefixIdx := strings.Index(line, LINE_PREFIX)
		splitInput := strings.Split(line[prefixIdx+1:], DELIMITER)
		scratchCards = append(scratchCards, NewScratchCard(
			list.Map[string, int](strings.Fields(splitInput[1]), mustBeInt),
			list.Map[string, int](strings.Fields(splitInput[0]), mustBeInt),
		))
	}

	return scratchCards
}

func NewScratchCard(numbers []int, winningNumbers []int) *ScratchCard {
	set := new(bit.Set)
	for _, winningNum := range winningNumbers {
		set.Add(winningNum)
	}

	return &ScratchCard{Numbers: numbers, WinningNumbers: set}
}

func (c *ScratchCard) NumberOfWins() int {
	numWinning := 0
	for _, num := range c.Numbers {
		if c.WinningNumbers.Contains(num) {
			numWinning++
		}
	}
	return numWinning
}

func mustBeInt(val string, _ int) int {
	return convert.MustBeInt(val)
}
