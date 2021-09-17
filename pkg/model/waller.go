package model

import "time"

// Transactions records wallet transactions
type Transaction struct {
	Datetime      time.Time `json:"datetime"`
	Amount        float64   `json:"amount"`
	StartDatetime time.Time `json:"startDatetime"`
	EndDatetime   time.Time `json:"endDatetime"`
}

// TransactionHistory records wallet transactions history
type TransactionHistory struct {
	Datetime time.Time `json:"datetime"`
	Amount   float64   `json:"amount"`
}
