package pqueuespan

// TopKQueue maintains a fixed-size queue of items
// with k highest priorities.
type TopKQueue struct {
	PQueueSpan
	k int
}

func NewTopKQueue(k int) *TopKQueue {
	return &TopKQueue{
		*NewPQueueSpan(MINPQ),
		k,
	}
}

// DryPush checks whether a Push with the given priority
// will result in a materialized insertion to the
// TopKQueue
func (pq *TopKQueue) DryPush(lowerBound, upperBound float64) bool {
	if pq.Size() < pq.k {
		return true
	}
	_, bottom := pq.Head()
	//if bottom < priority {
	if bottom.lowerBound == lowerBound {
		return bottom.upperBound < upperBound
		//return true
	}
	return bottom.lowerBound < lowerBound
	//return false
}

// Push pushes a new item to the TopKQueue, but does not
// actually insert the item into the queue unless its
// priority qualifies for the top-k
func (pq *TopKQueue) Push(value interface{}, lowerBound, upperBound float64) {
	if !pq.DryPush(lowerBound, upperBound) {
		return
	}
	if pq.Size() == pq.k {
		pq.Pop()
	}
	pq.PQueueSpan.Push(value, lowerBound, upperBound)
}

func (pq *TopKQueue) Descending() (values []interface{}, lowerBounds, upperBounds []float64) {
	values = make([]interface{}, pq.Size())
	lowerBounds = make([]float64, pq.Size())
	upperBounds = make([]float64, pq.Size())
	for i := len(values) - 1; i >= 0; i-- {
		v, p := pq.Pop()
		values[i] = v
		lowerBounds[i] = p.lowerBound
		upperBounds[i] = p.upperBound
	}
	return
}
