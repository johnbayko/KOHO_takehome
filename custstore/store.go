package custstore

import (
    "time"
)

type CustStore interface {
    BalanceAdd(
        id string, customerId string, loadAmountCents int64, time time.Time,
    ) error
}

