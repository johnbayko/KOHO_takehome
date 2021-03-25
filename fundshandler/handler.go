package fundshandler

import (
    //"fmt"  // debug
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

const (
    DAY = time.Hour * 24
    Week = DAY * 7
)
const (
    LOAD_AMOUNT_CENTS_LIMIT_DAY = 500000
    LOAD_AMOUNT_CENTS_LIMIT_WEEK = 2000000
    TRANSACTION_COUNT_LIMIT_DAY = 3
)

func (handler *FundsHandler) checkPerDayLimit(
    customerId string, loadAmountCents int64, transTime time.Time,
) (bool, error) {
    startAt :=
        time.Date(
            transTime.Year(),
            transTime.Month(),
            transTime.Day(),
            0,
            0,
            0,
            0,
            transTime.Location(),
        )
    endBefore := startAt.Add(DAY)

    totalLoadAmountCents, err :=
        handler.store.GetLoadAmountForPeriod(customerId, startAt, endBefore)
    if err != nil {
        // Can't ensure true, assume false.
        return false, err
    }
    newLoadAmountCents := totalLoadAmountCents + loadAmountCents

    isPerDayLimitOk := (newLoadAmountCents < LOAD_AMOUNT_CENTS_LIMIT_DAY)
    return isPerDayLimitOk, nil
}

func (handler *FundsHandler) Load(
    transId string,
    customerId string,
    loadAmountCents int64,
    transTime time.Time,
) (bool, error) {
    isPerDayLimitOk, err :=
        handler.checkPerDayLimit(customerId, loadAmountCents, transTime)
    if err != nil {
        return isPerDayLimitOk, err
    }
    if !isPerDayLimitOk {
        return false, nil
    }

    err = handler.store.BalanceAdd(
        transId, customerId, loadAmountCents, transTime)

    // Accepted if no rule checks fail, so is true even if err is not nil.
    return true, err
}
