package custstore

import (
    "time"
)

type CustStore interface {
    Open() error
    Close()

    BalanceAdd(
        id string, customerId string, loadAmountCents int64, time time.Time,
    ) error
    GetLoadAmountForPeriod(
        customerId string, startAt time.Time, endBefore time.Time,
    ) (int64, error)
}


type duplicateError struct {
}

func (err *duplicateError) Error() string {
    return "Duplicate transaction"
}

var DuplicateError = &duplicateError { }

