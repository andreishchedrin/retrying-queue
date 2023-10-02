package retrying_queue

import "github.com/pkg/errors"

var (
	ErrQueueOverflow = errors.New("queue capacity overflow")
	ErrQueueIsEmpty  = errors.New("queue is empty")
)

type node struct {
	value interface{}
	next  *node
}

type expect struct {
	value        interface{}
	alreadyRetry int
}

type Queue struct {
	head     *node
	tail     *node
	expect   *expect
	retries  int
	length   int
	capacity int
}

func NewQueue(retries int, capacity int) *Queue {
	if retries == 0 {
		retries = 1
	}

	return &Queue{
		retries:  retries,
		capacity: capacity,
	}
}

func (q *Queue) Length() int {
	return q.length
}

func (q *Queue) IsEmpty() bool {
	return q.length == 0
}

func (q *Queue) Success() {
	if q.expect != nil {
		q.resetExpect()
	}
}

func (q *Queue) Fail() {
	if q.expect == nil {
		return
	}

	if q.expect.alreadyRetry == q.retries {
		q.Success()
		return
	}

	q.expect.alreadyRetry++
}

func (q *Queue) Enqueue(value interface{}) error {
	if q.length == q.capacity && q.capacity != 0 {
		return ErrQueueOverflow
	}

	temp := &node{value: value}
	if q.head == nil {
		q.head = temp
		q.tail = temp
	} else {
		q.tail.next = temp
		q.tail = temp
	}
	q.length++

	return nil
}

func (q *Queue) Dequeue() (interface{}, error) {
	if q.expect != nil {
		return q.expect.value, nil
	}

	value, isEmpty := q.peek()
	if isEmpty {
		return nil, ErrQueueIsEmpty
	}

	q.setExpect(value)
	q.head = q.head.next
	q.length--
	return value, nil
}

func (q *Queue) peek() (interface{}, bool) {
	if q.IsEmpty() {
		return 0, true
	}

	return q.head.value, false
}

func (q *Queue) setExpect(value interface{}) {
	q.expect = &expect{
		value:        value,
		alreadyRetry: 1,
	}
}

func (q *Queue) resetExpect() {
	q.expect = nil
}
