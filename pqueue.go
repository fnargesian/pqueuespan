package pqueuespan

import (
	"fmt"
	"sync"
)

// PQType represents a priority queue ordering kind (see MAXPQ and MINPQ)
type PQType int

const (
	MAXPQ PQType = iota
	MINPQ
)

type item struct {
	value    interface{}
	priority span
}

type span struct {
	lowerBound float64
	upperBound float64
}

// PQueueSpan is a heap priority queue data structure implementation.
// It can be whether max or min ordered and it is synchronized
// and is safe for concurrent operations.
// The priority field in PQueueSpan is a range. Items in the queue are ordered lexicographically.
type PQueueSpan struct {
	sync.RWMutex
	items      []*item
	elemsCount int
	comparator func(span, span) bool
}

func newItem(value interface{}, lowerBound, upperBound float64) *item {
	priority := span{
		lowerBound: lowerBound,
		upperBound: upperBound,
	}
	return &item{
		value:    value,
		priority: priority,
	}
}

func (i *item) String() string {
	return fmt.Sprintf("<item value:%s priority:%v>", i.value, i.priority)
}

// NewPQueueSpan creates a new priority queue with the provided pqtype
// ordering type
func NewPQueueSpan(pqType PQType) *PQueueSpan {
	var cmp func(span, span) bool

	if pqType == MAXPQ {
		cmp = max
	} else {
		cmp = min
	}

	items := make([]*item, 1)
	items[0] = nil // Heap queue first element should always be nil

	return &PQueueSpan{
		items:      items,
		elemsCount: 0,
		comparator: cmp,
	}
}

// Push the value item into the priority queue with provided priority.
func (pq *PQueueSpan) Push(value interface{}, lowerBound, upperBound float64) {
	item := newItem(value, lowerBound, upperBound)

	pq.Lock()
	pq.items = append(pq.items, item)
	pq.elemsCount += 1
	pq.swim(pq.size())
	pq.Unlock()
}

// Pop and returns the highest/lowest priority item (depending on whether
// you're using a MINPQ or MAXPQ) from the priority queue
func (pq *PQueueSpan) Pop() (interface{}, span) {
	pq.Lock()
	defer pq.Unlock()

	if pq.size() < 1 {
		s := span{
			lowerBound: -1.0,
			upperBound: -1.0,
		}
		return nil, s
	}

	var max *item = pq.items[1]

	pq.exch(1, pq.size())
	pq.items = pq.items[0:pq.size()]
	pq.elemsCount -= 1
	pq.sink(1)

	return max.value, max.priority
}

// Head returns the highest/lowest priority item (depending on whether
// you're using a MINPQ or MAXPQ) from the priority queue
func (pq *PQueueSpan) Head() (interface{}, span) {
	pq.RLock()
	defer pq.RUnlock()

	if pq.size() < 1 {
		s := span{
			lowerBound: -1.0,
			upperBound: -1.0,
		}
		return nil, s
	}

	headValue := pq.items[1].value
	headPriority := pq.items[1].priority

	return headValue, headPriority
}

// Size returns the elements present in the priority queue count
func (pq *PQueueSpan) Size() int {
	pq.RLock()
	defer pq.RUnlock()
	return pq.size()
}

// Check queue is empty
func (pq *PQueueSpan) Empty() bool {
	pq.RLock()
	defer pq.RUnlock()
	return pq.size() == 0
}

func (pq *PQueueSpan) size() int {
	return pq.elemsCount
}

func max(x, y span) bool {
	//return i < j
	if y.lowerBound == x.lowerBound {
		return y.upperBound > x.upperBound
	}
	return y.lowerBound > x.lowerBound
}

func min(x, y span) bool {
	if y.lowerBound == x.lowerBound {
		return y.upperBound < x.upperBound
	}
	return y.lowerBound < x.lowerBound
	//return i > j
}

func (pq *PQueueSpan) less(i, j int) bool {
	return pq.comparator(pq.items[i].priority, pq.items[j].priority)
}

func (pq *PQueueSpan) exch(i, j int) {
	var tmpItem *item = pq.items[i]

	pq.items[i] = pq.items[j]
	pq.items[j] = tmpItem
}

func (pq *PQueueSpan) swim(k int) {
	for k > 1 && pq.less(k/2, k) {
		pq.exch(k/2, k)
		k = k / 2
	}

}

func (pq *PQueueSpan) sink(k int) {
	for 2*k <= pq.size() {
		var j int = 2 * k

		if j < pq.size() && pq.less(j, j+1) {
			j++
		}

		if !pq.less(k, j) {
			break
		}

		pq.exch(k, j)
		k = j
	}
}
