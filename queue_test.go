package retrying_queue

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type data struct {
	field int
}

type QueueTestSuite struct {
	suite.Suite
	queue    *Queue
	testData []data
}

func TestQueueTestSuite(t *testing.T) {
	suite.Run(t, &QueueTestSuite{})
}

func (q *QueueTestSuite) SetupTest() {
	q.queue = NewQueue(3, 3)
	q.testData = []data{
		{
			1,
		},
		{
			2,
		},
	}

	for _, data := range q.testData {
		_ = q.queue.Enqueue(data)
	}
}

func (q *QueueTestSuite) TestQueue_Success() {
	first, err := q.queue.Dequeue()
	q.NoError(err)
	q.Equal(q.testData[0].field, first.(data).field)
	//do something with value

	//retry
	q.queue.Fail()

	//try to get another element(not required)
	same, err := q.queue.Dequeue()
	q.NoError(err)
	q.Equal(q.testData[0].field, same.(data).field)

	//report about success
	q.queue.Success()

	//get next element
	second, err := q.queue.Dequeue()
	q.NoError(err)
	q.Equal(q.testData[1].field, second.(data).field)
	q.queue.Success()

	q.True(q.queue.IsEmpty())
	_, err = q.queue.Dequeue()
	q.ErrorIs(err, ErrQueueIsEmpty)
}

func (q *QueueTestSuite) TestQueue_Overflow() {
	_ = q.queue.Enqueue(data{field: 3})
	q.ErrorIs(q.queue.Enqueue(data{field: 4}), ErrQueueOverflow)
}

func (q *QueueTestSuite) TestQueue_AllRetries() {
	first, err := q.queue.Dequeue()
	q.NoError(err)
	q.Equal(q.testData[0].field, first.(data).field)

	q.queue.Fail()
	q.queue.Fail()
	same, err := q.queue.Dequeue()
	q.NoError(err)
	q.Equal(q.testData[0].field, same.(data).field)

	//retry limit
	q.queue.Fail()

	second, err := q.queue.Dequeue()
	q.NoError(err)
	q.Equal(q.testData[1].field, second.(data).field)
}

func (q *QueueTestSuite) TestQueue_DefaultRetriesAndCapacity() {
	q.queue = NewQueue(0, 0)
	for _, data := range q.testData {
		err := q.queue.Enqueue(data)
		q.NoError(err)
	}

	first, err := q.queue.Dequeue()
	q.NoError(err)
	q.Equal(q.testData[0].field, first.(data).field)

	//try to retry
	q.queue.Fail()

	second, err := q.queue.Dequeue()
	q.NoError(err)
	q.Equal(q.testData[1].field, second.(data).field)
}
