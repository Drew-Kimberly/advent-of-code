package day8_2023

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Drew-Kimberly/advent-of-code/src/common/go/fs"
	"github.com/Drew-Kimberly/advent-of-code/src/common/go/list"
	"github.com/Drew-Kimberly/advent-of-code/src/common/go/maths"
)

const (
	Right = "R"
	Left  = "L"
)

const START = "AAA"
const END = "ZZZ"

// Series of Right/Left instructions
type Instructions []string

type Node struct {
	Val   string
	Left  *Node
	Right *Node
}

type Path struct {
	Start    *Node
	End      *Node
	Curr     *Node
	NumSteps int
}

func Day8_2023() {
	inputPath, err := filepath.Abs("./2023/day8/input.txt")
	if err != nil {
		panic(err)
	}

	inputLines, err := fs.ExtractInputLines(inputPath)
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("Part 1 value: %d", PartOne(inputLines)))
	fmt.Println(fmt.Sprintf("Part 2 value: %d", PartTwo(inputLines)))
}

func PartOne(inputLines []string) int {
	instructions := parseInstructions(inputLines)
	path := parsePath(inputLines)

	path.WalkToEnd(instructions)

	return path.NumSteps
}

func PartTwo(inputLines []string) int {
	instructions := parseInstructions(inputLines)
	startingNodes := parseStartingNodes(inputLines)

	// For each starting node we'll keep track of the distance between end nodes
	// i.e. ABZ -> ABZ
	zLoops := make([][]int, len(startingNodes))
	instructionIdx := 0
	numSteps := 0
	allLoopsCollected := false

	updateZLoop := func(i int, steps int) {
		zLoops[i] = append(zLoops[i], steps)
	}

	for !allLoopsCollected {
		for i, node := range startingNodes {
			if instructions[instructionIdx] == Left {
				startingNodes[i] = node.Left
			} else {
				startingNodes[i] = node.Right
			}

			if startingNodes[i].Val[2:] == "Z" {
				updateZLoop(i, numSteps+1)
			}
		}

		if instructionIdx < len(instructions)-1 {
			instructionIdx++
		} else {
			instructionIdx = 0
		}

		numSteps++

		allLoopsCollected = len(list.Filter(zLoops, func(val []int, i int) bool {
			return len(val) < 2
		})) == 0
	}

	return maths.LCM(list.Map(zLoops, func(val []int, i int) int {
		return val[1] - val[0]
	})...)
}

func parseInstructions(inputLines []string) Instructions {
	return strings.Split(inputLines[0], "")
}

func parsePath(inputLines []string) *Path {
	valToNeighbors := make(map[string][]string)

	for _, line := range inputLines[2:] {
		valToNeighbors[line[0:3]] = []string{line[7:10], line[12:15]}
	}

	return NewPath(valToNeighbors)
}

func parseStartingNodes(inputLines []string) []*Node {
	valToNeighbors := make(map[string][]string)

	for _, line := range inputLines[2:] {
		valToNeighbors[line[0:3]] = []string{line[7:10], line[12:15]}
	}

	nodes := EnumerateAllNodes(valToNeighbors)
	var startingNodes []*Node
	for val, node := range nodes {
		if strings.HasSuffix(val, "A") {
			startingNodes = append(startingNodes, node)
		}
	}

	return startingNodes
}

func NewPath(valToNeighbors map[string][]string) *Path {
	nodes := EnumerateAllNodes(valToNeighbors)
	return &Path{Start: nodes[START], End: nodes[END], Curr: nodes[START], NumSteps: 0}
}

func EnumerateAllNodes(valToNeighbors map[string][]string) map[string]*Node {
	valToNode := make(map[string]*Node)
	var initNode func(nodeVal string, leftVal string, rightVal string) *Node

	initNode = func(nodeVal string, leftVal string, rightVal string) *Node {
		_, exists := valToNode[nodeVal]
		if !exists {
			valToNode[nodeVal] = NewNode(nodeVal, nil, nil)
			valToNode[nodeVal].Left = initNode(leftVal, valToNeighbors[leftVal][0], valToNeighbors[leftVal][1])
			valToNode[nodeVal].Right = initNode(rightVal, valToNeighbors[rightVal][0], valToNeighbors[rightVal][1])
		}

		return valToNode[nodeVal]
	}

	for val, neighborVals := range valToNeighbors {
		initNode(val, neighborVals[0], neighborVals[1])
	}

	return valToNode
}

func (p *Path) WalkToEnd(instructions Instructions) {
	i := 0

	for p.Curr.Val != END {
		instruction := instructions[i]

		if instruction == Left {
			p.Curr = p.Curr.Left
		} else if instruction == Right {
			p.Curr = p.Curr.Right
		}

		if i < len(instructions)-1 {
			i++
		} else {
			i = 0
		}

		p.NumSteps++
	}
}

func (p *Path) Visit(func(node *Node)) {

}

func NewNode(val string, left *Node, right *Node) *Node {
	node := &Node{Val: val, Left: left, Right: right}
	if left == nil {
		node.Left = node
	}
	if right == nil {
		node.Right = node
	}
	return node
}
