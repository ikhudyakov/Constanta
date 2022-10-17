package server

import (
	"encoding/json"
	"log"
	"net/http"
	"paymentService/internal/repo"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

//Создание платежа
func CreateTransaction(w http.ResponseWriter, r *http.Request, db repo.DBmanager) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	var err error
	var transaction repo.Transaction
	if err = json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		log.Println((err.Error()))
		w.Write([]byte(err.Error()))
		return
	}

	transaction.InitDate = time.Now().String()
	transaction.ModDate = transaction.InitDate
	transaction.Status = repo.NEW
	transaction.ID, err = db.SaveTransaction(&transaction)
	if err != nil {
		log.Println((err.Error()))
		w.Write([]byte(err.Error()))
		return
	}

	json.NewEncoder(w).Encode(transaction)
}

//Получение списка всех платежей пользователя по его id или email
func GetTransactions(w http.ResponseWriter, r *http.Request, db repo.DBmanager) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	var err error
	params := mux.Vars(r)
	var transactionsByUser []repo.Transaction

	parsId, err := strconv.ParseInt(params["idOrEmail"], 10, 64)
	if err != nil {
		if transactionsByUser, err = db.GetTransactionsByUserEmail(params["idOrEmail"]); err != nil {
			log.Println((err.Error()))
			w.Write([]byte(err.Error()))
			return
		}
		json.NewEncoder(w).Encode(transactionsByUser)
		return
	}

	if transactionsByUser, err = db.GetTransactionsByUserID(parsId); err != nil {
		log.Println((err.Error()))
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(transactionsByUser)
}

//Проверка статуса платежа по id
func GetStatus(w http.ResponseWriter, r *http.Request, db repo.DBmanager) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	var err error
	var t repo.Transaction
	params := mux.Vars(r)
	parsId, _ := strconv.ParseInt(params["id"], 10, 64)
	if t, err = db.GetTransactionByID(parsId); err != nil {
		log.Println((err.Error()))
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(t.Status))
}

//Отмена платежа по его id
func DeleteTransaction(w http.ResponseWriter, r *http.Request, db repo.DBmanager) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	parsId, _ := strconv.ParseInt(params["id"], 10, 64)
	result, err := db.DeleteTransactionByID(parsId)
	if err != nil {
		log.Println((err.Error()))
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(result))
}

//Изменение статуса платежа
func UpdateStatus(w http.ResponseWriter, r *http.Request, db repo.DBmanager) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	parsId, _ := strconv.ParseInt(params["id"], 10, 64)
	var transaction repo.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		log.Println((err.Error()))
		w.Write([]byte(err.Error()))
		return
	}

	result, err := db.UpdateTransactionStatusByID(parsId, transaction.Status)
	if err != nil {
		log.Println((err.Error()))
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(result))
}
