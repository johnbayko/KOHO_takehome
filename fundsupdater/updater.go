package main

import (
    "encoding/json"
    "fmt"
    "io"
    "os"
    "regexp"
    "strconv"
    "time"
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
    amountRe = regexp.MustCompile(`\$(\d+).(\d{2})`)
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
            fmt.Fprintf(os.Stderr, "Decode input: %v\n", decodeErr)
            return false, decodeErr
        }
    }
    t.Id = jt.Id
    t.Customer_id = jt.Customer_id

    // Decode string to cents `$([\d]+).([\d]{2})`
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

func update(inputFileName string, outputFileName string) error {
    inputFile, openInputErr := os.OpenFile(inputFileName, os.O_RDONLY, 0)
    if openInputErr != nil {
        fmt.Fprintf(os.Stderr, "Input file: %v\n", openInputErr)
        return openInputErr
    }
    defer inputFile.Close()

    outputFile, openOutputErr :=
        os.OpenFile(outputFileName, os.O_WRONLY | os.O_CREATE | os.O_EXCL, 0)
    if openOutputErr != nil {
        fmt.Fprintf(os.Stderr, "Output file: %v\n", openOutputErr)
        return openOutputErr
    }
    defer outputFile.Close()

    var t Transaction

    inputDecoder := json.NewDecoder(inputFile)
    for {
        isEnd, decodeErr := getTransaction(inputDecoder, &t)
        if isEnd {
            return nil
        }
        if decodeErr != nil {
            fmt.Fprintf(os.Stderr, "Decode input: %v\n", decodeErr)
            continue  // Try next transaction
        }
        fmt.Printf(
            "id %v customer_id %v load_amount %v time %v\n",
            t.Id,
            t.Customer_id,
            t.Load_amount_cents,
            t.Time)  // debug

        // Handle transaction and output error if any




    }
    // Won't actually get here.
    return nil
}

