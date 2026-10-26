package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

type Task struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	CreatedAt time.Time       `json:"created_at"`
}
type TaskQueue struct {
	pool *redis.Pool
}

func NewTaskQueue(pool *redis.Pool) *TaskQueue {
	return &TaskQueue{pool: pool}
}

func (q *TaskQueue) Enqueue(task Task) error {

	conn := q.pool.Get()
	defer conn.Close()

	data, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	_, err = conn.Do("LPUSH", "tasks:default", data)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}
