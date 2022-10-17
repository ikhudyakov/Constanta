package repo

type DBmanager interface {
	SaveTransaction(transaction *Transaction) (int64, error)
	GetTransactionsByUserID(userId int64) ([]Transaction, error)
	GetTransactionsByUserEmail(userEmail string) ([]Transaction, error)
	GetTransactionByID(id int64) (Transaction, error)
	UpdateTransactionStatusByID(id int64, status string) (string, error)
	DeleteTransactionByID(id int64) (string, error)
}
