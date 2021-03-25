package fundshandler

import (
    "time"

    "github.com/johnbayko/KOHO_takehome/custstore"
)

type FundsHandler struct {
    // customer store details
    store custstore.CustStore
}

func NewFundsHandler(cs custstore.CustStore) *FundsHandler {
    return &FundsHandler{
        store: cs,
    }
}

func (handler *FundsHandler) Load(
    transId string,
    customerId string,
    loadAmountCents int64,
    transTime time.Time,
) (bool, error) {
    err := handler.store.BalanceAdd(
        transId, customerId, loadAmountCents, transTime)
    // Accepted if no rule checks fail, so is true even if err is not nil.
    return true, err
}
