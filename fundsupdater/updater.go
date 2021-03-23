package main

import (
    "encoding/json"
    "fmt"
    "io"
    "os"
)

type Transaction struct {
    Id string `json:"id"`
    Customer_id string `json:"customer_id"`
    Load_amount string `json:"load_amount"`
    Time string `json:"time"`
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
    decode_err := input_decoder.Decode(&t)
    if decode_err != nil {
        if decode_err == io.EOF {
            return nil
        } else {
            fmt.Fprintf(os.Stderr, "Decode input: %v\n", decode_err)
            return decode_err
        }
    }
    fmt.Printf(
        "id %v customer_id %v load_amount %v time %v\n",
        t.Id,
        t.Customer_id,
        t.Load_amount,
        t.Time)  // debug




    return nil
}

