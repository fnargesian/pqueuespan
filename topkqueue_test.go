package pqueuespan

import "testing"

func Test_TopKQueue(t *testing.T) {
	sequence := []int{3, 1, 4, 5}
	// case 1
	//lbs := []int{1, 0, 3, 2}
	//ubs := []int{4, 4, 4, 4}
	// case 2
	//lbs := []int{1, 1, 1, 1}
	//ubs := []int{2, 6, 4, 5}
	// case 3
	//lbs := []int{3, 4, 1, 1}
	//ubs := []int{5, 6, 4, 6}
	// case 3
	sequence = []int{3, 1, 4, 5, 2}
	lbs := []int{3, 4, 3, 1, 5}
	ubs := []int{5, 6, 5, 6, 7}
	k := 3
	queue := NewTopKQueue(k)
	for i, v := range sequence {
		queue.Push(v, float64(lbs[i]), float64(ubs[i]))
	}
	for !queue.Empty() {
		v, priority := queue.Pop()
		t.Log(v, priority)
	}
}

func Test_TopKQueue_Descending(t *testing.T) {
	sequence := []int{3, 1, 4, 5, 2}
	lbs := []int{3, 4, 3, 1, 5}
	ubs := []int{5, 6, 5, 6, 7}
	k := 3
	queue := NewTopKQueue(k)
	for i, v := range sequence {
		queue.Push(v, float64(lbs[i]), float64(ubs[i]))
	}
	es, dlbs, _ := queue.Descending()
	prev := 1000.0
	for i, p := range dlbs {
		if p > prev {
			t.Error("Descending does not return the correct order")
		}
		prev = p
		t.Log(es[i])
	}
}
