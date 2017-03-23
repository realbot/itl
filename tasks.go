package itl

import (
	"fmt"
	"log"
	"time"

	"github.com/adjust/rmq"
)

const (
	unackedLimit = 1000
)

type Tasks struct {
	taskQueue rmq.Queue
}

func NewTasks(tasksTag, redisURL string) *Tasks {
	connection := rmq.OpenConnection("jobs", "tcp", redisURL, 0)
	taskQueue := connection.OpenQueue("tasks" + tasksTag)
	return &Tasks{
		taskQueue: taskQueue,
	}
}

func (t Tasks) EnqueueTask(req string) {
	t.taskQueue.Publish(req)
}

type PayloadProcessor func(string) error

func (t Tasks) StartConsumers(numConsumers int, processorFn PayloadProcessor) {
	t.taskQueue.StartConsuming(unackedLimit, 500*time.Millisecond)
	for i := 0; i < numConsumers; i++ {
		name := fmt.Sprintf("consumer%d", i)
		t.taskQueue.AddConsumer(name, t.newConsumer(i, processorFn))
	}
}

type taskConsumer struct {
	name      string
	processor PayloadProcessor
}

func (t Tasks) newConsumer(tag int, processor PayloadProcessor) *taskConsumer {
	return &taskConsumer{
		name:      fmt.Sprintf("consumer%d", tag),
		processor: processor,
	}
}

func (consumer *taskConsumer) Consume(delivery rmq.Delivery) {
	err := consumer.processor(delivery.Payload())
	if err != nil {
		log.Printf("%s: %v\n", consumer.name, err)
		delivery.Reject()
	} else {
		delivery.Ack()
	}
}
