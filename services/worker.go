package services

import (
	"context"
	"log"
	"time"
)

type Worker struct {
	queue *TaskQueue
}

func NewWorker(queue *TaskQueue) *Worker {
	return &Worker{queue: queue}
}

func (w *Worker) Run(ctx context.Context) {
	log.Println("воркер запущен, слушаю очередь tasks:default")

	for {
		select {
		case <-ctx.Done():
			log.Println("воркер остановлен")
		default:
		}
		task, err := w.queue.Dequeue(3 * time.Second)
		if err != nil {
			log.Printf("ошибка при получении задач %v", err)
			time.Sleep(time.Second)
			continue
		}
		if task == nil {
			continue
		}
	}
}
