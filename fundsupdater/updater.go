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
    amount_re = regexp.MustCompile(`\$(\d+).(\d{2})`)
)

/*
    Get the next transaction from the given json.Decoder, and fill the given
    Transaction fields.

    Returns boolean to indicate end of transactions, and error for current
    transaction.
 */
func get_transaction(input_decoder *json.Decoder, t *Transaction) (bool, error) {
    type JsonTransaction struct {
        Id string `json:"id"`
        Customer_id string `json:"customer_id"`
        Load_amount string `json:"load_amount"`
        Time time.Time `json:"time"`
    }
    var jt JsonTransaction

    decode_err := input_decoder.Decode(&jt)
    if decode_err != nil {
        if decode_err == io.EOF {
            return true, nil
        } else {
            fmt.Fprintf(os.Stderr, "Decode input: %v\n", decode_err)
            return false, decode_err
        }
    }
    t.Id = jt.Id
    t.Customer_id = jt.Customer_id

    // Decode string to cents `$([\d]+).([\d]{2})`
    amount_match := amount_re.FindStringSubmatch(jt.Load_amount)
    if len(amount_match) != 3 {
        return false, fmt.Errorf(
            "Amount not a valid \"$dollar.cents\" string: %v", jt.Load_amount)
    }

    // 2 digits for cents is about 7 bits, so allow no more than 64-7 = 57 bits
    // for dollars.
    amount_dollars, dollars_err := strconv.ParseInt(amount_match[1], 10, 57)
    if dollars_err != nil {
        return false, dollars_err
    }
    // The regex ensures this is 2 decimal digits, so parse will succeed,
    // error can be ignored.
    amount_cents, _ := strconv.ParseInt(amount_match[2], 10, 7)

    t.Load_amount_cents = (amount_dollars * 100) + amount_cents
    t.Time = jt.Time

    return false, nil
}

func update(input_file_name string, output_file_name string) error {
    input_file, open_input_err := os.OpenFile(input_file_name, os.O_RDONLY, 0)
    if open_input_err != nil {
        fmt.Fprintf(os.Stderr, "Input file: %v\n", open_input_err)
        return open_input_err
    }
    defer input_file.Close()

    output_file, open_output_err :=
        os.OpenFile(output_file_name, os.O_WRONLY | os.O_CREATE | os.O_EXCL, 0)
    if open_output_err != nil {
        fmt.Fprintf(os.Stderr, "Output file: %v\n", open_output_err)
        return open_output_err
    }
    defer output_file.Close()

    var t Transaction

    input_decoder := json.NewDecoder(input_file)
    for {
        is_end, decode_err := get_transaction(input_decoder, &t)
        if is_end {
            return nil
        }
        if decode_err != nil {
            fmt.Fprintf(os.Stderr, "Decode input: %v\n", decode_err)
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

