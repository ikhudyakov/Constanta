package repo

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
