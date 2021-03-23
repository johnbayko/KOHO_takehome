package main

import (
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
