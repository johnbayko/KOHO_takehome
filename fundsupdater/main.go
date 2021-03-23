package main

import (
    "fmt"
    "os"
)

const (
    INPUT_FILE_NAME = "input.txt"
    OUTPUT_FILE_NAME = "output.txt"
)

func main() {
    input_file_name := INPUT_FILE_NAME
    output_file_name := OUTPUT_FILE_NAME

    if len(os.Args) > 1 {
        input_file_name = os.Args[1]
    }
    if len(os.Args) > 2 {
        output_file_name = os.Args[2]
    }

    err := update(input_file_name, output_file_name)
    if err != nil {
        os.Exit(1)
    }

    os.Exit(0)
}

func update(input_file_name string, output_file_name string) error {
    input_file, err := os.OpenFile(input_file_name, os.O_RDONLY, 0)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Input file: %v\n", err)
        return err
    }
    defer input_file.Close()

    output_file, err :=
        os.OpenFile(output_file_name, os.O_WRONLY | os.O_CREATE | os.O_EXCL, 0)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Output file: %v\n", err)
        return err
    }
    defer output_file.Close()

    fmt.Println("Hello, world!")  // debug
    return nil
}
