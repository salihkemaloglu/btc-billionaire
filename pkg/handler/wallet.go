package handler

import (
	"fmt"
	"net/http"
	"sort"
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
		transaction.Amount = t.Amount
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
		ths := getTransactionsInBetween(t, transactions)
		ctx.JSON(http.StatusOK, ths)
		return
	}

	ctx.JSON(http.StatusInternalServerError, errors.New(fmt.Sprintf("no record for walletID: %s", walletID)))
}

// getTransactionsInBetween return transactions between start and end time
func getTransactionsInBetween(t model.Transaction, transactions []model.Transaction) []model.TransactionHistory {
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Datetime.Before(transactions[j].Datetime)
	})

	ths := []model.TransactionHistory{}
	year, month, day, hour := 0, 0, 0, 0
	var totalAmount float64 = 0
	for i := range transactions {
		if inTimeSpan(t.StartDatetime, t.EndDatetime, transactions[i].Datetime) {
			totalAmount += transactions[i].Amount
			if hour == 0 || transactions[i].Datetime.Year() != year || int(transactions[i].Datetime.Month()) != month ||
				transactions[i].Datetime.Day() != day || transactions[i].Datetime.Hour() != hour {
				year, month, day, hour = transactions[i].Datetime.Year(), int(transactions[i].Datetime.Month()), transactions[i].Datetime.Day(),
					transactions[i].Datetime.Hour()
				datetime := time.Date(transactions[i].Datetime.Year(), transactions[i].Datetime.Month(), transactions[i].Datetime.Day(),
					transactions[i].Datetime.Hour(), 0, 0, 0, transactions[i].Datetime.Location())
				ths = append(ths, model.TransactionHistory{Datetime: datetime, Amount: totalAmount})
			} else if transactions[i].Datetime.Year() == year || int(transactions[i].Datetime.Month()) == month ||
				transactions[i].Datetime.Day() == day || transactions[i].Datetime.Hour() == hour {
				setAmount(totalAmount, transactions[i], ths)
			}

		}
	}
	return ths
}

// setAmount sets amount of transaction that are same hour
func setAmount(totalAmount float64, transaction model.Transaction, ths []model.TransactionHistory) {
	for i := range ths {
		if ths[i].Datetime.Hour() == transaction.Datetime.Hour() {
			ths[i].Amount = totalAmount
		}
	}
}

// inTimeSpan checks time is between start and end
func inTimeSpan(start, end, check time.Time) bool {
	return (check.After(start) || check.Equal(start)) && (check.Before(end) || check.Equal(end))
}
