package models

type Transaction struct {
	ID          int `json:"id" db:"id"`
	UserId      int `json:"user_id" db:"user_id"`
	Time        int `json:"time" db:"time"`
	BalanceDiff int `json:"balance_diff" db:"balance_diff"`
}
