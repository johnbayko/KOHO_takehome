package fundshandler

import(
    "testing"
    "time"
)

var (
    // Let custStoreMock access the test instance.
    globalT *testing.T
    globalTestNum int
)

// Note - there are extra types because initialisation requires explicit type
// names even when struct literal is compatible.
type getLoadAmountForPeriodArgsMock struct {
    customerId string
    startAt time.Time
    endBefore time.Time
}
type getLoadAmountForPeriodReturnsMock struct {
    amount int64
    err error
}
type getLoadAmountForPeriodMock struct {
    args []getLoadAmountForPeriodArgsMock
    argsIdx int
    returns []getLoadAmountForPeriodReturnsMock
    returnsIdx int
}
type getNumForPeriodArgsMock struct {
    customerId string
    startAt time.Time
    endBefore time.Time
}
type getNumForPeriodReturnsMock struct {
    amount int64
    err error
}
type getNumForPeriodMock struct {
    args []getNumForPeriodArgsMock
    argsIdx int
    returns []getNumForPeriodReturnsMock
    returnsIdx int
}
type addTransactionArgsMock struct {
    id string
    customerId string
    loadAmountCents int64
    time time.Time
    accepted bool
}
type addTransactionReturnsMock struct {
    err error
}
type addTransactionMock struct {
    args []addTransactionArgsMock
    argsIdx int
    returns []addTransactionReturnsMock
    returnsIdx int
}
type balanceAddArgsMock struct {
    customerId string
    loadAmountCents int64
}
type balanceAddReturnsMock struct {
    err error
}
type balanceAddMock struct {
    args []balanceAddArgsMock
    argsIdx int
    returns []balanceAddReturnsMock
    returnsIdx int
}

type custStoreMock struct {
    getLoadAmountForPeriod getLoadAmountForPeriodMock
    getNumForPeriod getNumForPeriodMock
    addTransaction addTransactionMock
    balanceAdd balanceAddMock
}

// Not used for testing.
func (csm *custStoreMock) Open() error {
    return nil
}
// Not used for testing.
func (csm *custStoreMock) Close() {
}


func (csm *custStoreMock) GetLoadAmountForPeriod(
    customerId string, startAt time.Time, endBefore time.Time,
) (int64, error) {
    mock := csm.getLoadAmountForPeriod
    if mock.args[mock.argsIdx].customerId != customerId {
        globalT.Errorf("%v: GetLoadAmountForPeriod() customerId expected %v got %v",
            globalTestNum, mock.args[mock.argsIdx].customerId, customerId)
    }
    if !mock.args[mock.argsIdx].startAt.Equal(startAt) {
        globalT.Errorf("%v: GetLoadAmountForPeriod() startAt expected %v got %v",
            globalTestNum, mock.args[mock.argsIdx].startAt, startAt)
    }
    if !mock.args[mock.argsIdx].endBefore.Equal(endBefore) {
        globalT.Errorf("%v: GetLoadAmountForPeriod() endBefore expected %v got %v",
            globalTestNum, mock.args[mock.argsIdx].endBefore, endBefore)
    }
    mock.argsIdx = mock.argsIdx + 1

    // Need to use current returnsIdx, but increment before returning
    // (defer doesn't like doing that)
    returnsIdx := mock.returnsIdx
    mock.returnsIdx = mock.returnsIdx + 1

    return mock.returns[returnsIdx].amount, mock.returns[returnsIdx].err
}

func (csm *custStoreMock) GetNumForPeriod(
    customerId string, startAt time.Time, endBefore time.Time,
) (int64, error) {
    mock := csm.getNumForPeriod
    if mock.args[mock.argsIdx].customerId != customerId {
        globalT.Errorf("%v: GetNumForPeriod() customerId expected %v got %v",
            globalTestNum, mock.args[mock.argsIdx].customerId, customerId)
    }
    if !mock.args[mock.argsIdx].startAt.Equal(startAt) {
        globalT.Errorf("%v: GetNumForPeriod() startAt expected %v got %v",
            globalTestNum, mock.args[mock.argsIdx].startAt, startAt)
    }
    if !mock.args[mock.argsIdx].endBefore.Equal(endBefore) {
        globalT.Errorf("%v: GetNumForPeriod() endBefore expected %v got %v",
            globalTestNum, mock.args[mock.argsIdx].endBefore, endBefore)
    }
    // TODO Index gets incremented here, but reset to 0 next call for some
    // reason, making test fail incorrectly.
    mock.argsIdx = mock.argsIdx + 1

    // Need to use current returnsIdx, but increment before returning
    // (defer doesn't like doing that)
    returnsIdx := mock.returnsIdx
    mock.returnsIdx = mock.returnsIdx + 1

    return mock.returns[returnsIdx].amount, mock.returns[returnsIdx].err
}
func (csm *custStoreMock) AddTransaction(
    id string,
    customerId string,
    loadAmountCents int64,
    time time.Time,
    accepted bool,
) error {
    mock := csm.addTransaction
    if mock.args[mock.argsIdx].id != id {
        globalT.Errorf("%v: GetLoadAmountForPeriod() id expected %v got %v",
            globalTestNum, mock.args[mock.argsIdx].id, id)
    }
    if mock.args[mock.argsIdx].customerId != customerId {
        globalT.Errorf(
            "%v: GetLoadAmountForPeriod() customerId expected %v got %v",
            globalTestNum, mock.args[mock.argsIdx].customerId, customerId)
    }
    if mock.args[mock.argsIdx].loadAmountCents != loadAmountCents {
        globalT.Errorf(
            "%v: GetLoadAmountForPeriod() loadAmountCents expected %v got %v",
            globalTestNum,
            mock.args[mock.argsIdx].loadAmountCents,
            loadAmountCents,
        )
    }
    if !mock.args[mock.argsIdx].time.Equal(time) {
        globalT.Errorf("%v: GetLoadAmountForPeriod() time expected %v got %v",
            globalTestNum, mock.args[mock.argsIdx].time, time)
    }
    if mock.args[mock.argsIdx].accepted != accepted {
        globalT.Errorf(
            "%v: GetLoadAmountForPeriod() accepted expected %v got %v",
            globalTestNum, mock.args[mock.argsIdx].accepted, accepted)
    }
    mock.argsIdx = mock.argsIdx + 1

    // Need to use current returnsIdx, but increment before returning
    // (defer doesn't like doing that)
    returnsIdx := mock.returnsIdx
    mock.returnsIdx = mock.returnsIdx + 1

    return mock.returns[returnsIdx].err
}

func (csm *custStoreMock) BalanceAdd(
    customerId string, loadAmountCents int64,
) error {
    mock := csm.balanceAdd
    if mock.args[mock.argsIdx].customerId != customerId {
        globalT.Errorf(
            "%v: GetLoadAmountForPeriod() customerId expected %v got %v",
            globalTestNum, mock.args[mock.argsIdx].customerId, customerId)
    }
    if mock.args[mock.argsIdx].loadAmountCents != loadAmountCents {
        globalT.Errorf(
            "%v: GetLoadAmountForPeriod() loadAmountCents expected %v got %v",
            globalTestNum,
            mock.args[mock.argsIdx].loadAmountCents,
            loadAmountCents,
        )
    }
    mock.argsIdx = mock.argsIdx + 1

    // Need to use current returnsIdx, but increment before returning
    // (defer doesn't like doing that)
    returnsIdx := mock.returnsIdx
    mock.returnsIdx = mock.returnsIdx + 1

    return mock.returns[returnsIdx].err
}


func TestLoad(t *testing.T) {
    globalT = t

    type handlerArgs struct {
        transId string
        customerId string
        loadAmountCents int64
        transTime time.Time
    }
    type handlerReturns struct {
        b bool
        e error
    }
    var testList = []struct{
        args handlerArgs
        returns handlerReturns
        csm custStoreMock
    }{
        {
            handlerArgs {
                "123",
                "234",
                0,
                time.Date(2020, 01, 01, 01, 0, 0, 0, time.UTC),
            },
            handlerReturns {true, nil},
            custStoreMock {
                getLoadAmountForPeriodMock {
                    []getLoadAmountForPeriodArgsMock {
                        {
                            "234",
                            time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC),
                            time.Date(2020, 01, 02, 0, 0, 0, 0, time.UTC),
                        },
                        {
                            "234",
                            time.Date(2019, 12, 30, 0, 0, 0, 0, time.UTC),
                            time.Date(2020, 01, 05, 0, 0, 0, 0, time.UTC),
                        },
                    },
                    0, // argIdx
                    []getLoadAmountForPeriodReturnsMock {
                        {0, nil},
                        {0, nil},
                    },
                    0, // returnsIdx
                },
                getNumForPeriodMock {
                    []getNumForPeriodArgsMock {
                        {
                            "234",
                            time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC),
                            time.Date(2020, 01, 02, 0, 0, 0, 0, time.UTC),
                        },
                    },
                    0, // argIdx
                    []getNumForPeriodReturnsMock {
                        {0, nil},
                    },
                    0, // returnsIdx
                },
                addTransactionMock {
                    []addTransactionArgsMock {
                        {
                            "123",
                            "234",
                            0,
                            time.Date(2020, 01, 01, 01, 0, 0, 0, time.UTC),
                            true,
                        },
                    },
                    0, // argIdx
                    []addTransactionReturnsMock {
                        {nil},
                    },
                    0, // returnsIdx
                },
                balanceAddMock {
                    []balanceAddArgsMock {
                        {"234", 0},
                    },
                    0, // argIdx
                    []balanceAddReturnsMock {
                        {nil},
                    },
                    0, // returnsIdx
                },
            },
        },
    }
    globalTestNum := 0
    for _, test := range testList {
        csm := &test.csm
        handler := NewFundsHandler(csm)

        b, e := handler.Load(
                test.args.transId,
                test.args.customerId,
                test.args.loadAmountCents,
                test.args.transTime,
            )
        if test.returns.b != b {
            t.Errorf("%v: handler.Load() isAccepted expected %v got %v",
                globalTestNum, test.returns.b, b)
        }
        if test.returns.e != e {
            t.Errorf("%v: handler.Load() error expected %v got %v",
                globalTestNum, test.returns.e, e)
        }
        globalTestNum = globalTestNum + 1
    }
}
