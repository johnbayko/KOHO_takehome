package custstoresqlite

import (
    "time"

//    "github.com/mattn/go-sqlite3"
)

type CustStoreSqlite struct {
}

func NewCustStoreSqlite() *CustStoreSqlite {
    return &CustStoreSqlite {
    }
}

/*
    Update customer balance, save transaction if successful

    Returns error if fails.
 */
func (cs *CustStoreSqlite) BalanceAdd(
    id string, customerId string, loadAmountCents int64, time time.Time,
) error {
    return nil
}
