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
) error {
    err := handler.store.BalanceAdd(
        transId, customerId, loadAmountCents, transTime)
    return err
}
