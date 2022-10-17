package server

import (
	"net/http"
	c "paymentService/internal/config"
	"paymentService/internal/repo"

	"github.com/gorilla/mux"
)

func HandlersInit(dbmanager repo.DBmanager, conf *c.Config) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		CreateTransaction(w, r, dbmanager)
	}).Methods("POST")
	r.HandleFunc("/{idOrEmail}", func(w http.ResponseWriter, r *http.Request) {
		GetTransactions(w, r, dbmanager)
	}).Methods("GET")
	r.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		DeleteTransaction(w, r, dbmanager)
	}).Methods("DELETE")
	r.HandleFunc("/status/{id}", func(w http.ResponseWriter, r *http.Request) {
		UpdateStatus(w, r, dbmanager)
	}).Methods("PUT")
	r.HandleFunc("/status/{id}", func(w http.ResponseWriter, r *http.Request) {
		GetStatus(w, r, dbmanager)
	}).Methods("GET")

	return r
}
