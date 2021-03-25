package custstoresqlite

import (
    "database/sql"
    "fmt"  // debug
    "time"

    _ "github.com/mattn/go-sqlite3"

    "github.com/johnbayko/KOHO_takehome/custstore"
)

type CustStoreSqlite struct {
    db *sql.DB

    checkTransId *sql.Stmt
    checkCustId *sql.Stmt
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

    checkCustId, err := db.Prepare("select customer_id from customer where customer_id = ?")
    if err != nil {
        db.Close()
        return err
    }
    cs.checkCustId = checkCustId

    return nil
}

func (cs *CustStoreSqlite) Close() {
    cs.checkTransId.Close()
    cs.checkCustId.Close()
    cs.db.Close()
}

// Not exported

func (cs *CustStoreSqlite) isDuplicate(id string) (bool, error) {
    transIdRows, err := cs.checkTransId.Query(id)
    if err != nil {
        return false, err
    }
    defer transIdRows.Close()

    if transIdRows.Next() {
        // Transaciton ID already there, don't insert again (but not an error)
        // Need to add an indicator that it's not applied. Maybe a specific
        // error.
        return true, nil
    }
    return false, nil
}

func (cs *CustStoreSqlite) hasCustomer(customerId string) (bool, error) {
    custIdRows, err := cs.checkCustId.Query(customerId)
    if err != nil {
        return false, err
    }
    defer custIdRows.Close()

    if custIdRows.Next() {
        // Transaciton ID already there, don't insert again (but not an error)
        // Need to add an indicator that it's not applied. Maybe a specific
        // error.
        return true, nil
    }
    return false, nil
}

// Exported

/*
    Update apply a transaction to customer balance, save transaction record if
    successful. Reject duplicate transaction ids.

    Returns error if fails.
 */
func (cs *CustStoreSqlite) BalanceAdd(
    id string, customerId string, loadAmountCents int64, time time.Time,
) error {
    // Check transaction id
    isDuplicate, err := cs.isDuplicate(id)
    if err != nil {
        return err
    }
    if isDuplicate {
        return custstore.DuplicateError
    }

    // Update customers
    hasCustomer, err := cs.hasCustomer(id)
    if err != nil {
        return err
    }
    if !hasCustomer {
        // There is no cuseomer, need to create customer and account.
        fmt.Printf("No customer id %v\n", customerId)  // debug
    }


    // Update accounts

    // Add to transactions

    return nil
}
