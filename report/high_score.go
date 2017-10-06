package report

// An HighScoreItem is the aggregated result of a report.Card.
type HighScoreItem struct {
	Repo  string  `json:"repo"`
	Score float64 `json:"score"`
	Files int     `json:"files"`
}

// An HighScoreHeap is a min-heap of HighScoreItems.
type HighScoreHeap []HighScoreItem

func (h HighScoreHeap) Len() int           { return len(h) }
func (h HighScoreHeap) Less(i, j int) bool { return h[i].Score < h[j].Score }
func (h HighScoreHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

// Push onto the heap
func (h *HighScoreHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(HighScoreItem))
}

// Pop item off of the heap
func (h *HighScoreHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
