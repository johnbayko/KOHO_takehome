package main

import(
    "encoding/json"
    "strings"
    "testing"
    "time"
)

func TestGetTransaction(t *testing.T) {
    testString :=
`{"id":"15887","customer_id":"528","load_amount":"$3318.47","time":"2000-01-01T00:00:00Z"}
{"id":"15887","customer_id":"528","load_amount":"$3318.4","time":"2000-01-01T00:00:00Z"}
{"id":"15887","customer_id":"528","load_amount":"$.47","time":"2000-01-01T00:00:00Z"}
{"id":"15887","customer_id":"528","load_amount":"3318.47","time":"2000-01-01T00:00:00Z"}
{"id":"15887","customer_id":"528","load_amount":"3318","time":"2000-01-01T00:00:00Z"}
{"id":"15887","customer_id":"528","load_amount":"$3318.47","time":"20000-01-01T00:00:00Z"}`
    testDecoder := json.NewDecoder(strings.NewReader(testString))

    var expectList = []struct {
        isEnd bool
        has_err bool
        transaction Transaction
    }{
        // Okay
        {false, false,
            Transaction{
                "15887", "528", 331847,
                time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC) } },
        // Too few cents digits.
        {false, true, Transaction{"", "", 0, time.Time{}}},
        // No dollar digits.
        {false, true, Transaction{"", "", 0, time.Time{}}},
        // No dollar sign.
        {false, true, Transaction{"", "", 0, time.Time{}}},
        // No cents.
        {false, true, Transaction{"", "", 0, time.Time{}}},

        // Invalid date string
        {false, true, Transaction{"", "", 0, time.Time{}}},

        // End of input
        {true, false, Transaction{"", "", 0, time.Time{}}},
    }
    testNum := 0
    for _, expect := range expectList {
        var transaction Transaction
        isEnd, err := getTransaction(testDecoder, &transaction)
        if expect.isEnd != isEnd {
            t.Errorf("%v: isEnd expected %v got %v",
                testNum, expect.isEnd, isEnd)
        }
        if expect.has_err != (err != nil) {
            t.Errorf("%v: err expected %v got %v",
                testNum, expect.has_err, err != nil)
        }
        if (err == nil) && !isEnd {
            /// Check actual transaction for correctness.
            if expect.transaction.Id != transaction.Id {
                t.Errorf("%v: transaction.Id expected %v got %v",
                    testNum, expect.transaction.Id, transaction.Id)
            }
            if expect.transaction.Customer_id != transaction.Customer_id {
                t.Errorf("%v: transaction.Customer_id expected %v got %v",
                    testNum, expect.transaction.Customer_id, transaction.Customer_id)
            }
            if expect.transaction.Load_amount_cents != transaction.Load_amount_cents {
                t.Errorf("%v: transaction.Load_amount_cents expected %v got %v",
                    testNum, expect.transaction.Load_amount_cents, transaction.Load_amount_cents)
            }
            if !expect.transaction.Time.Equal(transaction.Time) {
                t.Errorf("%v: transaction.Time expected %v got %v",
                    testNum, expect.transaction.Time, transaction.Time)
            }
        }
        testNum = testNum + 1
    }
}
