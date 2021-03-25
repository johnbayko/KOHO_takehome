package custstoresqlite

import (
    "database/sql"
//    "fmt"  // debug
    "time"

    _ "github.com/mattn/go-sqlite3"

    "github.com/johnbayko/KOHO_takehome/custstore"
)

type CustStoreSqlite struct {
    db *sql.DB

    checkTransIdStmt *sql.Stmt
    checkCustIdStmt *sql.Stmt
    createCustomerStmt *sql.Stmt
    createAccountStmt *sql.Stmt
    updateAccountStmt *sql.Stmt
    addTransactionStmt *sql.Stmt

    loadAmountPerPeriodStmt *sql.Stmt
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

    checkTransIdStmt, err := db.Prepare("select id from transactions where id = ?")
    if err != nil {
        db.Close()
        return err
    }
    cs.checkTransIdStmt = checkTransIdStmt

    checkCustIdStmt, err := db.Prepare("select customer_id from customers where customer_id = ?")
    if err != nil {
        db.Close()
        return err
    }
    cs.checkCustIdStmt = checkCustIdStmt

    createCustomerStmt, err := db.Prepare("insert into customers (customer_id) values (?)")
    if err != nil {
        db.Close()
        return err
    }
    cs.createCustomerStmt = createCustomerStmt

    createAccountStmt, err := db.Prepare("insert into accounts (customer_id, balance) values (?, 0)")
    if err != nil {
        db.Close()
        return err
    }
    cs.createAccountStmt = createAccountStmt

    updateAccountStmt, err := db.Prepare("update accounts set balance = balance + ? where customer_id = ?")
    if err != nil {
        db.Close()
        return err
    }
    cs.updateAccountStmt = updateAccountStmt

    addTransactionStmt, err := db.Prepare("insert into transactions (id, customer_id, load_amount, time) values (?, ?, ?, ?)")
    if err != nil {
        db.Close()
        return err
    }
    cs.addTransactionStmt = addTransactionStmt

    loadAmountPerPeriodStmt, err := db.Prepare("select sum(load_amount) from transactions where customer_id = ? and time >= ? and time < ?")
    if err != nil {
        db.Close()
        return err
    }
    cs.loadAmountPerPeriodStmt = loadAmountPerPeriodStmt


    return nil
}

func (cs *CustStoreSqlite) Close() {
    cs.checkTransIdStmt.Close()
    cs.checkCustIdStmt.Close()
    cs.createCustomerStmt.Close()
    cs.createAccountStmt.Close()
    cs.updateAccountStmt.Close()
    cs.addTransactionStmt.Close()

    cs.loadAmountPerPeriodStmt.Close()

    cs.db.Close()
}

// Not exported

func (cs *CustStoreSqlite) isDuplicate(id string) (bool, error) {
    transIdRows, err := cs.checkTransIdStmt.Query(id)
    if err != nil {
        return false, err
    }
    defer transIdRows.Close()

    if transIdRows.Next() {
        return true, nil
    }
    return false, nil
}

func (cs *CustStoreSqlite) hasCustomer(customerId string) (bool, error) {
    custIdRows, err := cs.checkCustIdStmt.Query(customerId)
    if err != nil {
        return false, err
    }
    defer custIdRows.Close()

    if custIdRows.Next() {
        return true, nil
    }
    return false, nil
}

func (cs *CustStoreSqlite) createCustomerAndAccount(customerId string) error {
    _, err := cs.createCustomerStmt.Exec(customerId)
    if err != nil {
        return err
    }

    _, err = cs.createAccountStmt.Exec(customerId)
    return err
}

func (cs *CustStoreSqlite) updateAccount(
    loadAmountCents int64, customerId string,
) error {
    _, err := cs.updateAccountStmt.Exec(loadAmountCents, customerId)
    return err
}

func (cs *CustStoreSqlite) addTransaction(
    id string, customerId string, loadAmountCents int64, time time.Time,
) error {
    _, err := cs.addTransactionStmt.Exec(id, customerId, loadAmountCents, time)
    return err
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
    hasCustomer, err := cs.hasCustomer(customerId)
    if err != nil {
        return err
    }
    if !hasCustomer {
        // There is no cuseomer, need to create customer and account.
        err = cs.createCustomerAndAccount(customerId)
        if err != nil {
            return err
        }
    }

    // Update accounts
    err = cs.updateAccount(loadAmountCents, customerId)
    if err != nil {
        // Customer and account records will remain created.
        return err
    }

    // Add to transactions
    err = cs.addTransaction(id, customerId, loadAmountCents, time)
    if err != nil {
        // Customer and account records will remain created and updated,
        // no rollback.
        return err
    }

    return nil
}

func (cs *CustStoreSqlite) GetLoadAmountForPeriod(
    customerId string, startAt time.Time, endBefore time.Time,
) (int64, error) {
    loadAmountRows, err :=
        cs.loadAmountPerPeriodStmt.Query(customerId, startAt, endBefore)
    if err != nil {
        return 0, err
    }
    defer loadAmountRows.Close()

    // sum() will always return only one row on success.
    if !loadAmountRows.Next() {
        // Nothing? Shouldn't be possible, but assume nothing found.
        return 0, nil
    }
    var loadAmountCents int64
    loadAmountRows.Scan(&loadAmountCents)

    return loadAmountCents, nil
}
