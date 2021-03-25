package custstoresqlite

import (
    "database/sql"
    "time"

    _ "github.com/mattn/go-sqlite3"
)

type CustStoreSqlite struct {
    db *sql.DB
}

func NewCustStoreSqlite() *CustStoreSqlite {
    return &CustStoreSqlite { }
}

func (cs *CustStoreSqlite) Open() error {
    db, err := sql.Open("sqlite3", "cust_store.db")
    if err != nil {
        return err
    }
    cs.db = db
    return nil
}

func (cs *CustStoreSqlite) Close() {
    cs.db.Close()
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
