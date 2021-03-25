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
    WEEK = DAY * 7
)
const (
    LOAD_AMOUNT_CENTS_LIMIT_DAY = 500000
    LOAD_AMOUNT_CENTS_LIMIT_WEEK = 2000000
    TRANSACTION_COUNT_LIMIT_DAY = 3
)

func (handler *FundsHandler) checkPeriodLimit(
    customerId string,
    loadAmountCents int64,
    startAt time.Time,
    endBefore time.Time,
    loadAmountCentsLimit int64,
) (bool, error) {
    totalLoadAmountCents, err :=
        handler.store.GetLoadAmountForPeriod(customerId, startAt, endBefore)
    if err != nil {
        // Can't ensure true, assume false.
        return false, err
    }
    newLoadAmountCents := totalLoadAmountCents + loadAmountCents

    isPerDayLimitOk := (newLoadAmountCents < loadAmountCentsLimit)
    return isPerDayLimitOk, nil
}

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

    return handler.checkPeriodLimit(
            customerId,
            loadAmountCents,
            startAt,
            endBefore,
            LOAD_AMOUNT_CENTS_LIMIT_DAY,
        )
}

func (handler *FundsHandler) checkPerWeekLimit(
    customerId string, loadAmountCents int64, transTime time.Time,
) (bool, error) {
    transDay :=
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
    weekday := transTime.Weekday()

    var durationFromMonday time.Duration
    switch weekday {
    case time.Monday:
        durationFromMonday = DAY * 0
    case time.Tuesday:
        durationFromMonday = DAY * 1
    case time.Wednesday:
        durationFromMonday = DAY * 2
    case time.Thursday:
        durationFromMonday = DAY * 3
    case time.Friday:
        durationFromMonday = DAY * 4
    case time.Saturday:
        durationFromMonday = DAY * 5
    case time.Sunday:
        durationFromMonday = DAY * 6
    }

    startAt := transDay.Add(-durationFromMonday)
    endBefore := startAt.Add(WEEK)

    return handler.checkPeriodLimit(
            customerId,
            loadAmountCents,
            startAt,
            endBefore,
            LOAD_AMOUNT_CENTS_LIMIT_WEEK,
        )
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

    isPerWeekLimitOk, err :=
        handler.checkPerWeekLimit(customerId, loadAmountCents, transTime)
    if err != nil {
        return isPerWeekLimitOk, err
    }
    if !isPerWeekLimitOk {
        return false, nil
    }

    err = handler.store.BalanceAdd(
        transId, customerId, loadAmountCents, transTime)

    // Accepted if no rule checks fail, so is true even if err is not nil.
    return true, err
}
