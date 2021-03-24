package fundshandler

import (
    "time"
)

type FundsHandler struct {
    // customer store details
}

func NewFundsHandler(/*customer store*/) *FundsHandler {
    return &FundsHandler{
        // customer store
    }
}

func (handler *FundsHandler) Load(
    transId string,
    customerId string,
    loadAmountCents int64,
    transTime time.Time,
) bool {
    if len(transId) < 5 {
        return false
    }
    return true
}
