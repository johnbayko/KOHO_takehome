package fundshandler

import (
    "time"
)

func Load(
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
