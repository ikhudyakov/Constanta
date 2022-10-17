package pg

import (
	"database/sql"
	"paymentService/internal/repo"
	"time"
)

type PGmanager struct {
	DB *sql.DB
}

// Сохранение транзакции в базу данных
func (m *PGmanager) SaveTransaction(transaction *repo.Transaction) (int64, error) {

	var lastID int64
	err := m.DB.QueryRow(
		"INSERT INTO transactions (userid, useremail, amount, currency, initdate, moddate, status) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		transaction.UserId,
		transaction.UserEmail,
		transaction.Amount,
		transaction.Currency,
		transaction.InitDate,
		transaction.ModDate,
		transaction.Status).Scan(&lastID)

	if err != nil {
		return 0, err
	}

	return lastID, err
}

// Получение транзакций из базы данных по UserID
func (m *PGmanager) GetTransactionsByUserID(userId int64) ([]repo.Transaction, error) {
	transactions := []repo.Transaction{}

	rows, err := m.DB.Query("select * from Transactions where userid = $1", userId)
	if err != nil {
		return transactions, err
	}
	defer rows.Close()

	for rows.Next() {
		t := repo.Transaction{}
		err := rows.Scan(&t.ID, &t.UserId, &t.UserEmail, &t.Amount, &t.Currency, &t.InitDate, &t.ModDate, &t.Status)
		if err != nil {
			return transactions, err
		}
		transactions = append(transactions, t)
	}
	return transactions, err
}

// Получение транзакций из базы данных по UserEmail
func (m *PGmanager) GetTransactionsByUserEmail(userEmail string) ([]repo.Transaction, error) {
	transactions := []repo.Transaction{}
	rows, err := m.DB.Query("select * from Transactions where useremail = $1", userEmail)
	if err != nil {
		return transactions, err
	}
	defer rows.Close()

	for rows.Next() {
		t := repo.Transaction{}
		err := rows.Scan(&t.ID, &t.UserId, &t.UserEmail, &t.Amount, &t.Currency, &t.InitDate, &t.ModDate, &t.Status)
		if err != nil {
			return transactions, err
		}
		transactions = append(transactions, t)
	}
	return transactions, err
}

// Получение транзакции из базы данных по ее ID
func (m *PGmanager) GetTransactionByID(id int64) (repo.Transaction, error) {
	transaction := repo.Transaction{}
	rows, err := m.DB.Query("select * from Transactions where id = $1", id)
	if err != nil {
		return transaction, err
	}
	defer rows.Close()

	err = rows.Scan(
		&transaction.ID,
		&transaction.UserId,
		&transaction.UserEmail,
		&transaction.Amount,
		&transaction.Currency,
		&transaction.InitDate,
		&transaction.ModDate,
		&transaction.Status)

	if err != nil {
		return transaction, err
	}
	return transaction, sql.ErrConnDone
}

// Обновление статуса транзакции в базе данных по ее ID
func (m *PGmanager) UpdateTransactionStatusByID(id int64, status string) (string, error) {
	var statusFromBd, result string
	rows, err := m.DB.Query("select status from Transactions where id = $1", id)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	err = rows.Scan(&statusFromBd)
	if err != nil {
		return result, err
	}

	if statusFromBd != repo.FAIL && statusFromBd != repo.SUCCESS && statusFromBd != "" {
		_, err := m.DB.Exec("update Transactions set status = $1, moddate = $2 where id = $3", status, time.Now().String(), id)
		if err != nil {
			return result, err
		}
		result = "status updated successfully"
	}
	result = "status update error"
	return result, err
}

// Удаление транзакции из базе данных по ее ID
func (m *PGmanager) DeleteTransactionByID(id int64) (string, error) {
	var statusFromBd, result string
	rows, err := m.DB.Query("select status from Transactions where id = $1", id)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	err = rows.Scan(&statusFromBd)
	if err != nil {
		return result, err
	}

	if statusFromBd == "" {
		result = "transaction not found"
		return result, err
	}

	if statusFromBd != repo.FAIL && statusFromBd != repo.SUCCESS {
		r, err := m.DB.Exec("delete from Transactions where id = $1", id)
		if err != nil {
			return result, err
		}
		d, _ := r.RowsAffected()
		if d > 0 {
			result = "transaction deleted successfully"
			return result, err
		}
	}
	result = "transaction not deleted"
	return result, err
}
