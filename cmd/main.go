package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/transactions/{idOrEmail}", getTransactions).Methods("GET")
	r.HandleFunc("/transactions", createTransaction).Methods("POST")
	r.HandleFunc("/transactions/status/{id}", updateStatus).Methods("PUT")
	r.HandleFunc("/transactions/status/{id}", getStatus).Methods("GET")
	r.HandleFunc("/transactions/{id}", deleteTransaction).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8001", r))
}
