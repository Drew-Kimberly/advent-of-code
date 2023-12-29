package day17_2023

func NewVisitedNodes() *VisitedNodes {
	return &VisitedNodes{
		visited: make(map[*Node]map[Direction]map[int]bool),
	}
}

func (v *VisitedNodes) IsVisited(node *Node, direction Direction, numSteps int) bool {
	v.tryInitialize(node)
	return v.visited[node][direction][numSteps]
}

func (v *VisitedNodes) MarkVisited(node *Node, direction Direction, numSteps int) {
	v.tryInitialize(node)
	v.visited[node][direction][numSteps] = true
}

func (v *VisitedNodes) tryInitialize(node *Node) {
	_, initialized := v.visited[node]
	if !initialized {
		v.visited[node] = make(map[Direction]map[int]bool)
		v.visited[node][North] = make(map[int]bool)
		v.visited[node][East] = make(map[int]bool)
		v.visited[node][West] = make(map[int]bool)
		v.visited[node][South] = make(map[int]bool)
	}
}
