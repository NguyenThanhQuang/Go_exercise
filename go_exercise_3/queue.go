package main

type MyQueue struct {
	stackIn  []int
	stackOut []int
}

func Constructor() MyQueue {
	return MyQueue{
		stackIn:  make([]int, 0),
		stackOut: make([]int, 0),
	}
}

func (q *MyQueue) Push(x int) {
	q.stackIn = append(q.stackIn, x)
}

func (q *MyQueue) transferIfNeeded() {
	if len(q.stackOut) == 0 {
		for len(q.stackIn) > 0 {
			element := q.stackIn[len(q.stackIn)-1]
			q.stackIn = q.stackIn[:len(q.stackIn)-1]
			q.stackOut = append(q.stackOut, element)
		}
	}
}

func (q *MyQueue) Pop() int {
	q.transferIfNeeded()
	element := q.stackOut[len(q.stackOut)-1]
	q.stackOut = q.stackOut[:len(q.stackOut)-1]
	return element
}

func (q *MyQueue) Peek() int {
	q.transferIfNeeded()
	return q.stackOut[len(q.stackOut)-1]
}

func (q *MyQueue) Empty() bool {
	return len(q.stackIn) == 0 && len(q.stackOut) == 0
}
