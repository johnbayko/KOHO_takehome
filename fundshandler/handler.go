package fundshandler

import (
    "fmt"
    "os"
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
) bool {
    err := handler.store.BalanceAdd(
        transId, customerId, loadAmountCents, transTime)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Balance update: %v", err)
        return false
    }
    return true
}
