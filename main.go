package main

import (
	"assisko/handlers"
	"assisko/services"
	"assisko/store"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("ошибка .env")
	}
	port := os.Getenv("PORT")
	conn := os.Getenv("DATABASE_URL")
	redisAddr := os.Getenv("REDIS_ADDR")
	jwtSecret := os.Getenv("JWT_SECRET")

	initCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	st, err := store.NewPostgresStore(initCtx, conn)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer st.Close()

	svc := services.NewService(st)

	// настраиваю redis pool
	pool := &redis.Pool{
		MaxIdle:   10,
		MaxActive: 100,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisAddr)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
	defer pool.Close()
	// scs session Manager с Redis store

	session := scs.New()
	session.Store = redisstore.New(pool)
	session.Lifetime = 24 * time.Hour
	session.Cookie.Name = "assisko_session"
	session.Cookie.HttpOnly = true
	session.Cookie.Secure = false
	session.Cookie.SameSite = http.SameSiteLaxMode

	// Очередь задач

	queue := services.NewTaskQueue(pool)
	worker := services.NewWorker(queue)

	// передаем очередь в сервис, чтобы хендлеры могли кидать задачи

	svc.SetQueue(queue)
	// запуск фонового воркера

	workerCtx, workerCancel := context.WithCancel(context.Background())
	go worker.Run(workerCtx)

	h := handlers.NewHandler(svc, session, jwtSecret)
	// здесь в middleware кладем роутер
	handlerChain := session.LoadAndSave(h.Router())

	server := &http.Server{
		Addr:         port,
		Handler:      handlerChain,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error server %v", err)
		}
	}()
	log.Println("server started on :8080 (Postgres + Redis)")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("Получен сигнал остановки. Включается Graceful Shutdown")

	workerCancel()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server error %v", err)
	}

	log.Println("server correctly shutdown")

}
