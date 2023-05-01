/*
Copyright 2023 QKP Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package eventbus

import (
	"errors"
	"sync"
	"time"
)

// EventBus defines the methods for a basic eventbus data structure.
type EventBus interface {
	Add(item interface{}, expireAfter time.Duration) error
	Remove() (interface{}, error)
	Len() int
	IsEmpty() bool
	Get(index int) (interface{}, error)
	Peek() interface{}
	SetLimit(limit int)
}

// NewEventBus 实现EventBus Interface
func NewEventBus() EventBus {
	return &ArrayEventBus{}
}

// define eventItem structure.
type eventItem struct {
	value     interface{}
	timestamp time.Time
}

// ArrayEventBus represents an eventbus data structure that is implemented using an array.
type ArrayEventBus struct {
	items []eventItem
	lock  sync.Mutex
	limit int
	mutex sync.Mutex
}

// Add adds an element to the end of the eventbus.
func (q *ArrayEventBus) Add(item interface{}, expireAfter time.Duration) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.limit > 0 && len(q.items) >= q.limit {
		return errors.New("queue limit reached")
	}
	q.items = append(q.items, eventItem{
		value:     item,
		timestamp: time.Now().Add(expireAfter),
	})
	return nil
}

// Remove and returns the first element in the eventbus.
func (q *ArrayEventBus) Remove() (interface{}, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if len(q.items) == 0 {
		return nil, errors.New("queue is empty")
	}
	item := q.items[0].value
	q.items = q.items[1:]
	return item, nil
}

// Len returns the number of elements in the eventbus.
func (q *ArrayEventBus) Len() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return len(q.items)
}

// IsEmpty returns true if the queue is empty, false otherwise.
func (q *ArrayEventBus) IsEmpty() bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return len(q.items) == 0
}

// Peek returns the first element in the eventbus without removing it.
func (q *ArrayEventBus) Peek() interface{} {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if len(q.items) == 0 {
		return nil
	}
	return q.items[0]
}

func (q *ArrayEventBus) Get(index int) (interface{}, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if index < 0 || index >= len(q.items) {
		return nil, errors.New("index out of range")
	}
	return q.items[index].value, nil
}

func (q *ArrayEventBus) SetLimit(limit int) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.limit = limit
	if q.limit < 0 {
		q.limit = 0
	}

	if len(q.items) > q.limit {
		q.items = q.items[:q.limit]
	}

	now := time.Now()
	for len(q.items) > 0 && q.items[0].timestamp.Before(now) {
		q.items = q.items[1:]
	}
}
