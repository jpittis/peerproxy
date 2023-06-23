package peerproxy

import "time"

type DelayedEvent struct {
	event   *Event
	readyAt time.Time
}

type Queue struct {
	head  *DelayedEvent
	queue chan *DelayedEvent
	delay time.Duration
}

func NewQueue(delay time.Duration, size int) *Queue {
	return &Queue{
		delay: delay,
		queue: make(chan *DelayedEvent, size),
	}
}

func (q *Queue) Push(event *Event) {
	delayedEvent := &DelayedEvent{
		event:   event,
		readyAt: time.Now().Add(q.delay),
	}
	q.queue <- delayedEvent
}

func (q *Queue) Pop() *Event {
	for {
		if q.head == nil {
			q.head = <-q.queue
		}
		remaining := time.Until(q.head.readyAt)
		if remaining <= 0 {
			popped := q.head
			q.head = nil
			return popped.event
		} else {
			<-time.After(remaining)
		}
	}
}
