package day17_2023

// Min Priority Queue impl

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Cost < pq[j].Cost
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	no := x.(*QueueState)
	*pq = append(*pq, no)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	no := old[n-1]
	*pq = old[0 : n-1]
	return no
}
