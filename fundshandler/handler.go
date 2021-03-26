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

const (
    DAY = time.Hour * 24
    WEEK = DAY * 7
)
const (
    LOAD_AMOUNT_CENTS_LIMIT_DAY = 500000
    LOAD_AMOUNT_CENTS_LIMIT_WEEK = 2000000
    TRANSACTION_COUNT_LIMIT_DAY = 3
)

func (handler *FundsHandler) checkAmountPeriodLimit(
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

    isPerDayLimitOk := (newLoadAmountCents <= loadAmountCentsLimit)
    return isPerDayLimitOk, nil
}

func (handler *FundsHandler) checkAmountPerDayLimit(
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

    return handler.checkAmountPeriodLimit(
            customerId,
            loadAmountCents,
            startAt,
            endBefore,
            LOAD_AMOUNT_CENTS_LIMIT_DAY,
        )
}

func (handler *FundsHandler) checkAmountPerWeekLimit(
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

    return handler.checkAmountPeriodLimit(
            customerId,
            loadAmountCents,
            startAt,
            endBefore,
            LOAD_AMOUNT_CENTS_LIMIT_WEEK,
        )
}

func (handler *FundsHandler) checkNumPerDayLimit(
    customerId string, transTime time.Time,
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

    totalNum, err :=
        handler.store.GetNumForPeriod(customerId, startAt, endBefore)
    if err != nil {
        // Can't ensure true, assume false.
        return false, err
    }
    newNum := totalNum + 1

    isPerDayLimitOk := (newNum <= TRANSACTION_COUNT_LIMIT_DAY)
    return isPerDayLimitOk, nil
}

func (handler *FundsHandler) checkLimits(
    transId string,
    customerId string,
    loadAmountCents int64,
    transTime time.Time,
) (bool, error) {
    isAmountPerDayLimitOk, err :=
        handler.checkAmountPerDayLimit(customerId, loadAmountCents, transTime)
    if err != nil {
        return isAmountPerDayLimitOk, err
    }
    if !isAmountPerDayLimitOk {
        return false, nil
    }

    isAmountPerWeekLimitOk, err :=
        handler.checkAmountPerWeekLimit(customerId, loadAmountCents, transTime)
    if err != nil {
        return isAmountPerWeekLimitOk, err
    }
    if !isAmountPerWeekLimitOk {
        return false, nil
    }

    isNumPerDayLimitOk, err :=
        handler.checkNumPerDayLimit(customerId, transTime)
    if err != nil {
        return isNumPerDayLimitOk, err
    }
    if !isNumPerDayLimitOk {
        return false, nil
    }
    return true, nil
}

func (handler *FundsHandler) Load(
    transId string,
    customerId string,
    loadAmountCents int64,
    transTime time.Time,
) (bool, error) {
    isOk, err :=
        handler.checkLimits(transId, customerId, loadAmountCents, transTime)
    if err != nil {
        return isOk, err
    }
    if !isOk {
        err = handler.store.AddTransaction(
            transId, customerId, loadAmountCents, transTime, false)
        return false, err
    }

    err = handler.store.BalanceAdd(
        transId, customerId, loadAmountCents, transTime)

    // Accepted if no rule checks fail, so is true even if err is not nil.
    return true, err
}
