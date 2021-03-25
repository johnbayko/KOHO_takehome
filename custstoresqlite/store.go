package custstoresqlite

import (
    "database/sql"
    "fmt"  // debug
    "time"

    _ "github.com/mattn/go-sqlite3"
)

type CustStoreSqlite struct {
    db *sql.DB

    checkTransId *sql.Stmt
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

    checkTransId, err := db.Prepare("select id from transactions where id = ?")
    if err != nil {
        db.Close()
        return err
    }
    cs.checkTransId = checkTransId

    return nil
}

func (cs *CustStoreSqlite) Close() {
    cs.checkTransId.Close()
    cs.db.Close()
}

/*
    Update apply a transaction to customer balance, save transaction record if
    successful. Reject duplicate transaction ids.

    Returns error if fails.
 */
func (cs *CustStoreSqlite) BalanceAdd(
    id string, customerId string, loadAmountCents int64, time time.Time,
) error {
    // Check transaction id
    transIdRows, err := cs.checkTransId.Query(id)
    if err != nil {
        return err
    }
    if transIdRows.Next() {
        // Transaciton ID already there, don't insert again (but not an error)
        return nil
    }
    // Update customers
    // Update accounts
    // Add to transactions
    return nil
}
