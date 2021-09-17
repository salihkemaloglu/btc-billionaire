package handler

import (
	"encoding/json"

	"github.com/patrickmn/go-cache"
	"github.com/salihkemaloglu/btc-billionaire/pkg/model"
)

// static user wallet id
var walletID = "ca7bff33-fd23-4fbe-8b35-aea32fa85346"

// Handler is endpoints
type Handler struct {
	cache    *cache.Cache
	WalletID string
}

// NewHandler example
func NewHandler(c *cache.Cache) *Handler {
	return &Handler{cache: c, WalletID: walletID}
}

// DefaultSet sets default value to the cache
func DefaultSet(c *cache.Cache) {
	str := `
	{
		"datetime": "2021-09-16T18:25:05+07:00",
		"amount": 1000
	}`
	var t model.Transaction
	if err := json.Unmarshal([]byte(str), &t); err != nil {
		panic(err)
	}

	c.Set(walletID, []model.Transaction{t}, cache.NoExpiration)
}
