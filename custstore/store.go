package custstore

import (
    "time"
)

type CustStore interface {
    Open() error
    Close()

    AddTransaction(
        id string,
        customerId string,
        loadAmountCents int64,
        time time.Time,
        accepted bool,
    ) error

    BalanceAdd(customerId string, loadAmountCents int64) error

    GetLoadAmountForPeriod(
        customerId string, startAt time.Time, endBefore time.Time,
    ) (int64, error)

    GetNumForPeriod(
        customerId string, startAt time.Time, endBefore time.Time,
    ) (int64, error)
}


type duplicateError struct {
}

func (err *duplicateError) Error() string {
    return "Duplicate transaction"
}

var DuplicateError = &duplicateError { }

