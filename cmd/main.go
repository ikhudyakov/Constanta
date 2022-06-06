package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type Transaction struct {
	ID       string  `json:"id"`
	User     User    `json:"user"`
	Amount   float32 `json:"amount"`
	Currency string  `json:"currency"`
	InitDate string  `json:"initdate"`
	ModDate  string  `json:"moddate"`
	Status   string  `json:"status"`
}

const (
	NEW     = "NEW"
	SUCCESS = "SUCCESS"
	FAIL    = "FAIL"
	ERROR   = "ERROR"
)

var transactions []Transaction

//Создание платежа
func createTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var transaction Transaction
	_ = json.NewDecoder(r.Body).Decode(&transaction)

	transaction.ID = strconv.Itoa(rand.Intn(1000000))
	transaction.InitDate = time.Now().String()
	transaction.ModDate = time.Now().String()
	transaction.Status = NEW

	transactions = append(transactions, transaction)
	json.NewEncoder(w).Encode(transaction)
}

//Получение списка всех платежей пользователя по его id или email
func getTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var transactionsByUser []Transaction
	for _, item := range transactions {

		if item.User.ID == params["idOrEmail"] || item.User.Email == params["idOrEmail"] {
			transactionsByUser = append(transactionsByUser, item)
		}
	}
	json.NewEncoder(w).Encode(transactionsByUser)
}

//Проверка статуса платежа по id
func getStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var status string
	for _, item := range transactions {

		if item.ID == params["id"] {
			status = item.Status
		}
	}
	w.Write([]byte(status))
}

//Отмена платежа по его id
func deleteTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range transactions {
		if item.ID == params["id"] {
			if item.Status == FAIL || item.Status == SUCCESS {
				w.WriteHeader(http.StatusLocked)
				w.Write([]byte("Error! Transaction status: " + item.Status))
				return
			}
			transactions = append(transactions[:index], transactions[index+1:]...)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Transaction deleted"))
			return
		}
	}
	w.Write([]byte("Transaction not found"))
}

//Изменение статуса платежа
func updateStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range transactions {
		if item.ID == params["id"] {
			var transaction Transaction
			_ = json.NewDecoder(r.Body).Decode(&transaction)
			if transactions[index].Status == FAIL || transactions[index].Status == SUCCESS {
				w.WriteHeader(http.StatusLocked)
				w.Write([]byte("Error! Transaction status: " + transactions[index].Status))
				return
			}
			transactions[index].Status = transaction.Status
			transactions[index].ModDate = time.Now().String()
			json.NewEncoder(w).Encode(transactions[index])
			return
		}
	}
	w.Write([]byte("Transaction not found"))
}

func main() {
	TempTransaction("1", "1", "no-reply@test.ru", 152.86, "RUB", NEW)
	TempTransaction("2", "1", "no-reply@test.ru", 5474.14, "RUB", NEW)
	TempTransaction("3", "1", "no-reply@test.ru", 5622.52, "RUB", FAIL)
	TempTransaction("4", "1", "no-reply@test.ru", 985.25, "EUR", SUCCESS)
	TempTransaction("5", "1", "no-reply@test.ru", 551.47, "EUR", ERROR)

	r := mux.NewRouter()

	r.HandleFunc("/transactions/{idOrEmail}", getTransactions).Methods("GET")
	r.HandleFunc("/transactions", createTransaction).Methods("POST")
	r.HandleFunc("/transactions/status/{id}", updateStatus).Methods("PUT")
	r.HandleFunc("/transactions/status/{id}", getStatus).Methods("GET")
	r.HandleFunc("/transactions/{id}", deleteTransaction).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8001", r))
}

func TempTransaction(
	id string,
	userId string,
	userEmail string,
	amount float32,
	currency string,
	status string) {
	transactions = append(transactions, Transaction{
		ID:       id,
		User:     User{ID: userId, Email: userEmail},
		Amount:   amount,
		Currency: currency,
		InitDate: time.Now().String(),
		ModDate:  time.Now().String(),
		Status:   status,
	})
}
