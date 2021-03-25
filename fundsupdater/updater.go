package main

import (
    "encoding/json"
    "fmt"
    "io"
    "os"
    "regexp"
    "strconv"
    "time"

    "github.com/johnbayko/KOHO_takehome/fundshandler"
)

type Transaction struct {
    Id string
    Customer_id string

    // Input amount is no more than 8 characters, or 6 digits.
    // n digits take roughly n * log2(10) = 6 * 3.3 bits = 20,
    // so 64 is plenty in this case.
    Load_amount_cents int64

    Time time.Time
}

var (
    amountRe = regexp.MustCompile(`\$(\d+)\.(\d{2})`)
)

/*
    Get the next transaction from the given json.Decoder, and fill the given
    Transaction fields.

    Returns boolean to indicate end of transactions, and error for current
    transaction.
 */
func getTransaction(inputDecoder *json.Decoder, t *Transaction) (bool, error) {
    type JsonTransaction struct {
        Id string `json:"id"`
        Customer_id string `json:"customer_id"`
        Load_amount string `json:"load_amount"`
        Time time.Time `json:"time"`
    }
    var jt JsonTransaction

    decodeErr := inputDecoder.Decode(&jt)
    if decodeErr != nil {
        if decodeErr == io.EOF {
            return true, nil
        } else {
            return false, decodeErr
        }
    }
    t.Id = jt.Id
    t.Customer_id = jt.Customer_id

    // Decode string to cents.
    amountMatch := amountRe.FindStringSubmatch(jt.Load_amount)
    if len(amountMatch) != 3 {
        return false, fmt.Errorf(
            "Amount not a valid \"$dollar.cents\" string: %v", jt.Load_amount)
    }

    // 2 digits for cents is about 7 bits, so allow no more than 64-7 = 57 bits
    // for dollars.
    amountDollars, dollarsErr := strconv.ParseInt(amountMatch[1], 10, 57)
    if dollarsErr != nil {
        return false, dollarsErr
    }
    // The regex ensures this is 2 decimal digits, so parse will succeed,
    // error can be ignored.
    amountCents, _ := strconv.ParseInt(amountMatch[2], 10, 7)

    t.Load_amount_cents = (amountDollars * 100) + amountCents
    t.Time = jt.Time

    return false, nil
}

/*
    Write a transaction acceptance record to the given json.Encoder.

    Return error if there's a problem.
 */
func putAcceptance(
    outputEncoder *json.Encoder, t *Transaction, isAccepted bool,
) error {
    type JsonAcceptance struct {
        Id string `json:"id"`
        Customer_id string `json:"customer_id"`
        Accepted bool `json:"accepted"`
    }
    var ja JsonAcceptance = JsonAcceptance{t.Id, t.Customer_id, isAccepted}

    encodeErr := outputEncoder.Encode(&ja)
    if encodeErr != nil {
        return encodeErr
    }
    return nil
}

func update(
    inputFileName string,
    outputFileName string,
    handler *fundshandler.FundsHandler,
) error {
    inputFile, openInputErr := os.OpenFile(inputFileName, os.O_RDONLY, 0)
    if openInputErr != nil {
        fmt.Fprintf(os.Stderr, "Input file: %v\n", openInputErr)
        return openInputErr
    }
    defer inputFile.Close()

    inputDecoder := json.NewDecoder(inputFile)

    outputFile, openOutputErr :=
        os.OpenFile(outputFileName, os.O_WRONLY | os.O_CREATE | os.O_EXCL, 0666)
    if openOutputErr != nil {
        fmt.Fprintf(os.Stderr, "Output file: %v\n", openOutputErr)
        return openOutputErr
    }
    defer outputFile.Close()

    outputEncoder := json.NewEncoder(outputFile)

    var t Transaction
    for {
        isEnd, getTransErr := getTransaction(inputDecoder, &t)
        if isEnd {
            return nil
        }
        if getTransErr != nil {
            fmt.Fprintf(os.Stderr, "Get transaction: %v\n", getTransErr)
            continue  // Try next transaction
        }
        isAccepted :=
            handler.Load(t.Id, t.Customer_id, t.Load_amount_cents, t.Time)

        putAcceptErr := putAcceptance(outputEncoder, &t, isAccepted)
        if putAcceptErr != nil {
            fmt.Fprintf(os.Stderr, "Put acceptance record: %v\n", putAcceptErr)
        }
    }
    // Won't actually get here.
    return nil
}

