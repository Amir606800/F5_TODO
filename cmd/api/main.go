package main

import (
	"F5/internals/handler"
	"F5/internals/repository"
	"F5/internals/service"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error when loading .env: ", err)
	}

	ctx := context.Background()

	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	db := os.Getenv("POSTGRES_DB")

	connString := fmt.Sprintf(
		"postgres://%s:%s@localhost:5433/%s?sslmode=disable",
		user,
		pass,
		db,
	)

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	err = pool.Ping(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("successfully connected to db")

	repo := repository.NewRepository(pool)
	svc := service.NewService(repo)
	hand := handler.NewHandler(svc)

	mux := http.NewServeMux()

	// TODOS
	mux.HandleFunc("GET /tasks", hand.GetTasks)
	mux.HandleFunc("POST /tasks", hand.CreateTask)
	mux.HandleFunc("GET /tasks/{id}", hand.GetTask)
	mux.HandleFunc("DELETE /tasks/{id}", hand.DeleteTask)
	//mux.HandleFunc("PATCH /tasks/{id}", hand)
	//mux.HandleFunc("GET /tasks?status=done", hand.GetTasks)

	fmt.Println("Server starting on :8082")
	log.Fatal(http.ListenAndServe(os.Getenv("API_URL"), mux))
}
