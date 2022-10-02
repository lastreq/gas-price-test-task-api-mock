package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/", SendGasHistory).Methods("GET")

	srv := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Server Started")

	<-done
	log.Print("Server Stopped")
}

func SendGasHistory(w http.ResponseWriter, r *http.Request) {
	buffer, err := os.ReadFile("gas_price.json")
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(buffer)
	if err != nil {
		log.Println(err)
	}
}
