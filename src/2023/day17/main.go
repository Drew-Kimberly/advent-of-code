package day17_2023

import (
	"container/heap"
	"fmt"
	"path/filepath"

	"github.com/Drew-Kimberly/advent-of-code/src/common/go/fs"
)

func Day17_2023() {
	inputPath, err := filepath.Abs("./2023/day17/input.txt")
	if err != nil {
		panic(err)
	}

	inputLines, err := fs.ExtractInputLines(inputPath)
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("Part 1 value: %d", FindMinHeatLoss(NewCityMap(inputLines))))
}

func FindMinHeatLoss(cityMap CityMap) int {
	minHeatLoss := FindPath(cityMap[0][0], cityMap[len(cityMap)-1][len(cityMap[0])-1], cityMap)

	return minHeatLoss
}

func FindPath(from, to *Node, cityMap CityMap) int {
	queue := &PriorityQueue{}
	heap.Init(queue)
	heap.Push(queue, &QueueState{Node: from, Cost: 0, NumSteps: 0})

	visited := NewVisitedNodes()

	for queue.Len() > 0 {
		state := heap.Pop(queue).(*QueueState)

		if state.Node == to {
			// Found a path to the goal node.
			// Since priority queue is sorted by min cost
			// we guarantee this is the min heatloss and can just return.
			return state.Cost
		}

		if visited.IsVisited(state.Node, state.Direction, state.NumSteps) {
			continue
		}

		visited.MarkVisited(state.Node, state.Direction, state.NumSteps)

		for _, neighbor := range NeighborStates(state, cityMap) {
			neighbor.Cost = state.Cost + neighbor.Node.HeatLoss
			heap.Push(queue, neighbor)
		}
	}

	panic("No path found!")
}

func NeighborStates(state *QueueState, cityMap CityMap) []*QueueState {
	var neighbors []*QueueState
	curr := state.Node

	// To the north
	cannotMoveNorth := curr.Coord.y == 0 || state.Direction == South || (state.NumSteps == 3 && state.Direction == North)
	if !cannotMoveNorth {
		nextState := &QueueState{
			Node:      cityMap[curr.Coord.y-1][curr.Coord.x],
			Direction: North,
			NumSteps:  1,
		}
		if state.Direction == North {
			nextState.NumSteps = state.NumSteps + 1
		}

		neighbors = append(neighbors, nextState)
	}

	// To the south
	cannotMoveSouth := curr.Coord.y == len(cityMap)-1 || state.Direction == North || (state.NumSteps == 3 && state.Direction == South)
	if !cannotMoveSouth {
		nextState := &QueueState{
			Node:      cityMap[curr.Coord.y+1][curr.Coord.x],
			Direction: South,
			NumSteps:  1,
		}
		if state.Direction == South {
			nextState.NumSteps = state.NumSteps + 1
		}

		neighbors = append(neighbors, nextState)
	}

	// To the east
	cannotMoveEast := curr.Coord.x == len(cityMap[0])-1 || state.Direction == West || (state.NumSteps == 3 && state.Direction == East)
	if !cannotMoveEast {
		nextState := &QueueState{
			Node:      cityMap[curr.Coord.y][curr.Coord.x+1],
			Direction: East,
			NumSteps:  1,
		}
		if state.Direction == East {
			nextState.NumSteps = state.NumSteps + 1
		}

		neighbors = append(neighbors, nextState)
	}

	// To the west
	cannotMoveWest := curr.Coord.x == 0 || state.Direction == East || (state.NumSteps == 3 && state.Direction == West)
	if !cannotMoveWest {
		nextState := &QueueState{
			Node:      cityMap[curr.Coord.y][curr.Coord.x-1],
			Direction: West,
			NumSteps:  1,
		}
		if state.Direction == West {
			nextState.NumSteps = state.NumSteps + 1
		}

		neighbors = append(neighbors, nextState)
	}

	return neighbors
}
