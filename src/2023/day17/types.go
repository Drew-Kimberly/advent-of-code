package day17_2023

type CityMap = [][]*Node

type Coordinate struct {
	x int
	y int
}

type Node struct {
	Coord    *Coordinate
	HeatLoss int
}

// For a given node tracks whether it has been explored,
// accounting for direction and num steps in that direction
// to accommodate the max moves constraint.
type VisitedNodes struct {
	visited map[*Node]map[Direction]map[int]bool
}

type Direction int

const (
	North Direction = iota
	East
	South
	West
)

type QueueState struct {
	// Next cell on the map to explore
	Node *Node
	// Total heatloss for the given path
	Cost int
	// Direction of this cell from the previous cell
	Direction Direction
	// Number of steps we've made in the current direction
	NumSteps int
}

type PriorityQueue []*QueueState
