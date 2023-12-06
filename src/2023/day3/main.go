package day3_2023

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Drew-Kimberly/advent-of-code/src/common/go/convert"
	"github.com/Drew-Kimberly/advent-of-code/src/common/go/fs"
)

func Day3_2023() {
	inputPath, err := filepath.Abs("./2023/day3/input.txt")
	if err != nil {
		panic(err)
	}

	inputLines, err := fs.ExtractInputLines(inputPath)
	if err != nil {
		panic(err)
	}

	schematic := NewSchematic(inputLines)

	fmt.Println(fmt.Sprintf("Part 1 value: %d", sumPartNumbers(schematic)))
	fmt.Println(fmt.Sprintf("Part 2 value: %d", sumGearRatios(schematic)))
}

func sumPartNumbers(schematic *Schematic) int {
	sum := 0
	for _, partNum := range schematic.PartNumbers() {
		sum += partNum.Val
	}
	return sum
}

func sumGearRatios(schematic *Schematic) int {
	sum := 0
	for _, gear := range schematic.Gears() {
		sum += (gear.PartNumbers[0].Val * gear.PartNumbers[1].Val)
	}
	return sum
}

type Schematic struct {
	matrix [][]string
}

type SchematicCharacter[T string | int] struct {
	Row      int
	ColStart int
	ColEnd   int
	Val      T
}

type PartNumber SchematicCharacter[int]

type Gear struct {
	Char        SchematicCharacter[string]
	PartNumbers []*PartNumber
}

func NewSchematic(inputLines []string) *Schematic {
	matrix := make([][]string, len(inputLines))

	for i, line := range inputLines {
		for _, char := range line {
			matrix[i] = append(matrix[i], string(char))
		}
	}

	s := Schematic{matrix: matrix}
	return &s
}

func (s *Schematic) PartNumbers() []*PartNumber {
	var partNumbers []*PartNumber
	for i, row := range s.matrix {
		j := 0
		var numStr string

		for j < len(row) {
			char := row[j]
			if isIntChar(char) {
				numStr += char
			} else {
				if len(numStr) > 0 && s.isPartNumber(convert.MustBeInt(numStr), i, j-len(numStr)) {
					partNumber := PartNumber{Row: i, ColStart: j - len(numStr), ColEnd: j - 1, Val: convert.MustBeInt(numStr)}
					partNumbers = append(partNumbers, &partNumber)
				}

				numStr = ""
			}

			j++
		}

		if len(numStr) > 0 && s.isPartNumber(convert.MustBeInt(numStr), i, j-len(numStr)) {
			partNumber := PartNumber{Row: i, ColStart: j - len(numStr), ColEnd: j - 1, Val: convert.MustBeInt(numStr)}
			partNumbers = append(partNumbers, &partNumber)
		}

		numStr = ""
	}

	return partNumbers
}

func (s *Schematic) Gears() []*Gear {
	var gears []*Gear
	gearToAdjacentPartNum := make(map[string]map[string]*PartNumber)
	partNumbers := s.PartNumbers()

	for i := range partNumbers {
		j := i

		for j < len(partNumbers) {
			partNumber := partNumbers[j]
			for _, char := range s.adjacentChars(partNumber.Row, partNumber.ColStart, partNumber.ColEnd) {
				if !s.isGearSymbol(char.Val) {
					continue
				}

				gearKey := fmt.Sprintf("%d-%d", char.Row, char.ColStart)
				partNumKey := fmt.Sprintf("%d-%d", partNumber.Row, partNumber.ColStart)

				_, exists := gearToAdjacentPartNum[gearKey]
				if !exists {
					gearToAdjacentPartNum[gearKey] = make(map[string]*PartNumber)
				}
				gearToAdjacentPartNum[gearKey][partNumKey] = partNumber
			}

			j++
		}

	}

	for gearKey, val := range gearToAdjacentPartNum {
		// Gears must be adjacent to exactly 2 part numbers.
		if len(val) != 2 {
			continue
		}

		gearKeyArr := strings.Split(gearKey, "-")
		gear := Gear{
			Char: SchematicCharacter[string]{
				Row:      convert.MustBeInt(gearKeyArr[0]),
				ColStart: convert.MustBeInt(gearKeyArr[1]),
				ColEnd:   convert.MustBeInt(gearKeyArr[1]),
			},
			PartNumbers: []*PartNumber{},
		}

		for _, partNum := range val {
			gear.PartNumbers = append(gear.PartNumbers, partNum)
		}

		gears = append(gears, &gear)
	}

	return gears
}

func (s *Schematic) isPartNumber(num int, rowIdx int, colStartIdx int) bool {
	numDigits := len(strconv.Itoa(num))

	for _, adjacentChar := range s.adjacentChars(rowIdx, colStartIdx, colStartIdx+numDigits-1) {
		if s.isSymbol(adjacentChar.Val) {
			return true
		}
	}

	return false
}

func (s *Schematic) adjacentChars(row int, colStart int, colEnd int) []*SchematicCharacter[string] {
	var chars []*SchematicCharacter[string]

	if colStart > 0 {
		char := SchematicCharacter[string]{Row: row, ColStart: colStart - 1, ColEnd: colStart - 1, Val: s.matrix[row][colStart-1]}
		chars = append(chars, &char)
	}

	if colEnd < len(s.matrix[row])-1 {
		char := SchematicCharacter[string]{Row: row, ColStart: colEnd + 1, ColEnd: colEnd + 1, Val: s.matrix[row][colEnd+1]}
		chars = append(chars, &char)
	}

	if row > 0 {
		for col := colStart - 1; col <= colEnd+1; col++ {
			if col >= 0 && col < len(s.matrix[row-1]) {
				char := SchematicCharacter[string]{Row: row - 1, ColStart: col, ColEnd: col, Val: s.matrix[row-1][col]}
				chars = append(chars, &char)
			}
		}
	}

	if row < len(s.matrix)-1 {
		for col := colStart - 1; col <= colEnd+1; col++ {
			if col >= 0 && col < len(s.matrix[row+1]) {
				char := SchematicCharacter[string]{Row: row + 1, ColStart: col, ColEnd: col, Val: s.matrix[row+1][col]}
				chars = append(chars, &char)
			}
		}
	}

	return chars
}

func (s *Schematic) isSymbol(char string) bool {
	return char != "." && !isIntChar(char)
}

func (s *Schematic) isGearSymbol(char string) bool {
	return char == "*"
}

func charInPartNumber(char *SchematicCharacter[string], partNumbers []*PartNumber) bool {
	if !isIntChar(char.Val) {
		return false
	}

	for _, partNumber := range partNumbers {
		if char.Row == partNumber.Row && char.ColStart >= partNumber.ColStart && char.ColEnd <= partNumber.ColEnd {
			return true
		}
	}

	return false
}

func isIntChar(char string) bool {
	_, err := strconv.Atoi(char)
	return err == nil
}
