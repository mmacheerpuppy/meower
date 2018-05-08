package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tinrab/meower/meow-service/db"
	"github.com/tinrab/meower/mq"
)

func newRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/meows", createMeowHandler).
		Methods("POST").
		Queries("body", "{body}")
	return
}

func main() {
	repo, err := db.NewPostgres("postgres://meower:123456@postgres/meower?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	db.SetRepository(repo)
	defer repo.Close()

	queue, err := mq.NewRabbitMessageQueue("amqp://guest:guest@rabbitmq:5672")
	if err != nil {
		log.Fatal(err)
	}
	mq.SetMessageQueue(queue)
	defer queue.Close()

	router := newRouter()
	if err := http.ListenAndServe(":3000", router); err != nil {
		log.Fatal(err)
	}
}