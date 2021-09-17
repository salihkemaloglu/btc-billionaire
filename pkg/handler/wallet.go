package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"github.com/salihkemaloglu/btc-billionaire/pkg/httputil"
	"github.com/salihkemaloglu/btc-billionaire/pkg/model"
	"github.com/salihkemaloglu/btc-billionaire/pkg/model/validate"
)

// InsertTransaction inserts transactions to the wallet
func (h *Handler) InsertTransaction(ctx *gin.Context) {
	var t model.Transaction
	if err := ctx.ShouldBindJSON(&t); err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	if err := validate.Transaction(t); err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	if o, found := h.cache.Get(walletID); found {
		transactions, ok := o.([]model.Transaction)
		if !ok {
			httputil.NewError(ctx, http.StatusInternalServerError, errors.New("couldn't parse object"))
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
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	if o, found := h.cache.Get(walletID); found {
		transactions, ok := o.([]model.Transaction)
		if !ok {
			httputil.NewError(ctx, http.StatusInternalServerError, errors.New("couldn't parse object"))
			return
		}
		th := getTransactionsInBetween(t, transactions)
		ctx.JSON(http.StatusOK, th)
		return
	}

	ctx.JSON(http.StatusInternalServerError, errors.New(fmt.Sprintf("no record for walletID: %s", walletID)))
}

// getTransactionsInBetween return transactions between start and end time
func getTransactionsInBetween(t model.Transaction, transactions []model.Transaction) []model.TransactionHistory {
	th := []model.TransactionHistory{}
	hour := 0
	for _, transaction := range transactions {
		if inTimeSpan(t.StartDatetime, t.EndDatetime, transaction.Datetime) {
			if hour == 0 || transaction.Datetime.Hour() != hour {
				hour = transaction.Datetime.Hour()
				transaction.Datetime = time.Date(transaction.Datetime.Year(), transaction.Datetime.Month(), transaction.Datetime.Day(),
					transaction.Datetime.Hour(), 0, 0, 0, transaction.Datetime.Location())
			}
			th = append(th, model.TransactionHistory{Datetime: transaction.Datetime, Amount: transaction.Amount})
		}
	}
	return th
}

// inTimeSpan checks time is between start and end
func inTimeSpan(start, end, check time.Time) bool {
	return (check.After(start) || check.Equal(start)) && (check.Before(end) || check.Equal(end))
}
