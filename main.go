package main

import (
	"assisko/store"
	"context"
	"log"
	"os"
	"time"
)

func main() {
	conn := os.Getenv("DATABASE_URL")

	initCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	st, err := store.NewPostgresStore(initCtx, conn)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer st.Close()

	svc := services.NewService(st)

}
