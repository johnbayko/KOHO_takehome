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
}


type duplicateError struct {
}

func (err *duplicateError) Error() string {
    return "Duplicate transaction"
}

var DuplicateError = &duplicateError { }

