package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"github.com/salihkemaloglu/btc-billionaire/model"
)

// InsertTransaction inserts transactions to the wallet
func (h *Handler) InsertTransaction(ctx *gin.Context) {
	var t model.Transaction
	if err := ctx.ShouldBindJSON(&t); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	if o, found := h.cache.Get(walletID); found {
		transactions, ok := o.([]model.Transaction)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, errors.New("couldn't parse object"))
			return
		}

		transaction := transactions[len(transactions)-1]
		transaction.Datetime = t.Datetime
		transaction.Amount += t.Amount
		transactions = append(transactions, transaction)
		h.cache.Set(walletID, transactions, cache.NoExpiration)
	} else {
		h.cache.Set(walletID, []model.Transaction{t}, cache.NoExpiration)
	}

	ctx.JSON(http.StatusOK, "OK")
}

// GetTransactions returns wallet's transaction history
func (h *Handler) GetTransactions(ctx *gin.Context) {
	var t model.Transaction
	if err := ctx.ShouldBindJSON(&t); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	if o, found := h.cache.Get(walletID); found {
		transactions, ok := o.([]model.Transaction)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, errors.New("couldn't parse object"))
			return
		}

		th := []model.TransactionHistory{}
		hour := 0
		for _, transaction := range transactions {
			if inTimeSpan(t.StartDatetime, t.EndDatetime, transaction.Datetime) {
				if hour == 0 || transaction.Datetime.Hour() != hour {
					hour = transaction.Datetime.Hour()
					transaction.Datetime = time.Date(transaction.Datetime.Year(), transaction.Datetime.Month(), transaction.Datetime.Day(),
						transaction.Datetime.Hour(), 0, 0, transaction.Datetime.Nanosecond(), transaction.Datetime.Location())
				}
				th = append(th, model.TransactionHistory{Datetime: transaction.Datetime, Amount: transaction.Amount})
			}
		}

		ctx.JSON(http.StatusOK, th)
		return
	}

	ctx.JSON(http.StatusInternalServerError, errors.New(fmt.Sprintf("no record for walletID: %s", walletID)))
}

// inTimeSpan checks time is between start and end
func inTimeSpan(start, end, check time.Time) bool {
	return (check.After(start) || check.Equal(start)) && (check.Before(end) || check.Equal(end))
}
