package queue

import (
	"sync"
	"testing"
)

func TestQueue(t *testing.T) {
	q := New[int]()

	if !q.IsEmpty() {
		t.Errorf("Expected empty queue")
	}

	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	if q.Len() != 3 {
		t.Errorf("Expected length 3, got %d", q.Len())
	}

	v, ok := q.Peek()
	if !ok || v != 1 {
		t.Errorf("Peek expected 1, got %v", v)
	}

	v, ok = q.Dequeue()
	if !ok || v != 1 {
		t.Errorf("Dequeue expected 1, got %v", v)
	}

	v, ok = q.Dequeue()
	if !ok || v != 2 {
		t.Errorf("Dequeue expected 2, got %v", v)
	}

	q.Enqueue(4)

	v, ok = q.Dequeue()
	if !ok || v != 3 {
		t.Errorf("Dequeue expected 3, got %v", v)
	}

	v, ok = q.Dequeue()
	if !ok || v != 4 {
		t.Errorf("Dequeue expected 4, got %v", v)
	}

	if !q.IsEmpty() {
		t.Errorf("Expected empty queue")
	}
}

func TestQueueCapacity(t *testing.T) {
	q := NewWithCapacity[int](10)
	if q.Len() != 0 {
		t.Errorf("Expected length 0")
	}
}

func TestQueueClear(t *testing.T) {
	q := New[int]()
	q.Enqueue(1)
	q.Enqueue(2)
	q.Clear()

	if q.Len() != 0 || !q.IsEmpty() {
		t.Errorf("Expected empty queue after clear")
	}

	v, ok := q.Dequeue()
	if ok {
		t.Errorf("Expected false on empty dequeue, got %v", v)
	}
}

func TestQueueConcurrency(t *testing.T) {
	q := New[int]()
	var wg sync.WaitGroup
	numOps := 1000
	numGophers := 100

	for i := 0; i < numGophers; i++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			for j := 0; j < numOps; j++ {
				q.Enqueue(base + j)
			}
		}(i * 10000)
	}

	wg.Wait()

	if q.Len() != numOps*numGophers {
		t.Errorf("Expected %d elements, got %d", numOps*numGophers, q.Len())
	}

	for i := 0; i < numGophers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numOps; j++ {
				q.Dequeue()
			}
		}()
	}

	wg.Wait()

	if q.Len() != 0 {
		t.Errorf("Expected 0 elements, got %d", q.Len())
	}
}

func TestQueueWrapAround(t *testing.T) {
	q := NewWithCapacity[int](3)
	q.Enqueue(1)
	q.Enqueue(2)
	q.Dequeue()  // head=1, tail=2, count=1
	q.Enqueue(3) // head=1, tail=0, count=2
	q.Enqueue(4) // head=1, tail=1, count=3, FULL

	// This should trigger resize
	q.Enqueue(5) // head=0, tail=4, count=4

	expected := []int{2, 3, 4, 5}
	for _, exp := range expected {
		v, ok := q.Dequeue()
		if !ok || v != exp {
			t.Errorf("Expected %d, got %v", exp, v)
		}
	}
}
