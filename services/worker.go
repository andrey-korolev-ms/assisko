package services

import (
	"context"
	"encoding/json"
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
		w.process(task)
	}
}
func (w *Worker) process(task *Task) {
	log.Printf("воркер выполняет задачу %v %v", task.ID, task.Type)
	switch task.Type {
	case "send_welcome_email":
		w.handleWelcomeEmail(task)
	default:
		log.Printf("неизвестный тип задачи: %s", task.Type)
	}
}
func (w *Worker) handleWelcomeEmail(task *Task) {
	var payload struct {
		UserID int    `json:"user_id"`
		Email  string `json:"email"`
	}
	if err := json.Unmarshal(task.Payload, &payload); err != nil {
		log.Printf("не удалось распарсить payload: %v ", err)
		return
	}
	log.Printf("Отправка сообщения на %s (user_id=%d)", payload.Email, payload.UserID)
	time.Sleep(15 * time.Second)
	log.Printf("Сообщение для пользователя %d отправлено", payload.UserID)
}
