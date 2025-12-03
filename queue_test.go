package main

import "testing"

func TestQueue(t *testing.T) {
	q := NewQueue[string]()
	q.Enqueue("item1", "item2")
	q.Enqueue("item3")

	item1, _ := q.Dequeue()
	if item1 != "item1" {
		t.Error("expected item1")
	}
	q.Dequeue()
	item3, _ := q.Dequeue()
	if item3 != "item3" {
		t.Error("expected item3")
	}
	_, err := q.Dequeue()
	if err == nil {
		t.Error("expected error")
	}
}
