package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Transaction struct {
	ID        int64   `json:"id"`
	UserId    int64   `json:"userId"`
	UserEmail string  `json:"userEmail"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	InitDate  string  `json:"initdate"`
	ModDate   string  `json:"moddate"`
	Status    string  `json:"status"`
}

const (
	NEW     = "NEW"
	SUCCESS = "SUCCESS"
	FAIL    = "FAIL"
	ERROR   = "ERROR"
)

//Создание платежа
func createTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var transaction Transaction
	_ = json.NewDecoder(r.Body).Decode(&transaction)
	transaction.InitDate = time.Now().String()
	transaction.ModDate = transaction.InitDate
	transaction.Status = NEW
	transaction.ID = SaveToDB(&transaction)

	json.NewEncoder(w).Encode(transaction)
}

//Получение списка всех платежей пользователя по его id или email
func getTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var transactionsByUser []Transaction

	parsId, err := strconv.ParseInt(params["idOrEmail"], 10, 64)
	if err != nil {
		transactionsByUser = GetFromDBByUserEmail(params["idOrEmail"])
		json.NewEncoder(w).Encode(transactionsByUser)
		return
	}

	transactionsByUser = GetFromDBByUserID(parsId)
	json.NewEncoder(w).Encode(transactionsByUser)
}

//Проверка статуса платежа по id
func getStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	parsId, _ := strconv.ParseInt(params["id"], 10, 64)
	status := GetFromDBByID(parsId).Status
	w.Write([]byte(status))
}

//Отмена платежа по его id
func deleteTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	parsId, _ := strconv.ParseInt(params["id"], 10, 64)

	w.Write([]byte(DeleteFromDBByID(parsId)))
}

//Изменение статуса платежа
func updateStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	parsId, _ := strconv.ParseInt(params["id"], 10, 64)
	var transaction Transaction
	_ = json.NewDecoder(r.Body).Decode(&transaction)

	w.Write([]byte(UpdateStatusDBByID(parsId, transaction.Status)))
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/transactions/{idOrEmail}", getTransactions).Methods("GET")
	r.HandleFunc("/transactions", createTransaction).Methods("POST")
	r.HandleFunc("/transactions/status/{id}", updateStatus).Methods("PUT")
	r.HandleFunc("/transactions/status/{id}", getStatus).Methods("GET")
	r.HandleFunc("/transactions/{id}", deleteTransaction).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8001", r))
}
