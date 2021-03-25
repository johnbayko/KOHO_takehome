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

