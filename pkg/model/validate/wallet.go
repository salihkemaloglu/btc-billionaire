package validate

import (
	"errors"

	"github.com/salihkemaloglu/btc-billionaire/pkg/model"
)

func Transaction(t model.Transaction) error {
	switch {
	case t.Amount == 0 || t.Amount < 0:
		return errors.New("transaction amount can't equal or smaller then 0")
	default:
		return nil
	}
}
