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

// кладем задачу в Redis
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

// удаляем задачу из Redis
func (q *TaskQueue) Dequeue(timeout time.Duration) (*Task, error) {
	conn := q.pool.Get()
	defer conn.Close()

	reply, err := redis.Values(conn.Do("BRPOP", "tasks:default", int(timeout.Seconds())))
	if err == redis.ErrNil {
		// таймаут, значит задач нет
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	data, err := redis.Bytes(reply[1], nil)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	var task Task
	if err := json.Unmarshal(data, &task); err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return &task, nil
}
