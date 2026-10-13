package main

import (
	"assisko/handlers"
	"assisko/services"
	"assisko/store"
	"context"
	"log"
	"net/http"
	"os"
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

	h := handlers.NewHandler(svc, session, jwtSecret)
	// здесь в middleware кладем роутер
	handlerChain := session.LoadAndSave(h.Router())

	server := &http.Server{
		Addr:         port,
		Handler:      handlerChain,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Println("server started on :8080 (Postgres + Redis)")
	log.Fatal(server.ListenAndServe())

}
