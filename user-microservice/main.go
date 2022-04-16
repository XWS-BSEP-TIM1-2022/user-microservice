package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	router := mux.NewRouter()
	router.StrictSlash(true)
	server, err := NewUserServer()
	if err != nil {
		log.Fatal(err.Error())
	}

	defer server.CloseTracer()
	defer server.CloseDB()

	router.HandleFunc("/users/{id}", server.getUserHandler).Methods("GET")
	router.HandleFunc("/users", server.getAllUsersHandler).Methods("GET")
	router.HandleFunc("/users", server.createUserHandler).Methods("POST")
	router.HandleFunc("/users/{id}", server.updateUserHandler).Methods("PUT")
	router.HandleFunc("/users/{id}", server.deleteUserHandler).Methods("DELETE")

	// start userServer
	srv := &http.Server{Addr: "0.0.0.0:8000", Handler: router}
	go func() {
		log.Println("userServer starting")
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	<-quit

	log.Println("service shutting down ...")

	// gracefully stop userServer
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("userServer stopped")
}
