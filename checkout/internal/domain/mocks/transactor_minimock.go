package mocks

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

//go:generate minimock -i gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/domain.Transactor -o ./mocks\transactor_minimock.go -n TransactorMock

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
)

// TransactorMock implements domain.Transactor
type TransactorMock struct {
	t minimock.Tester

	funcRunReadCommitted          func(ctx context.Context, f func(ctxTX context.Context) error) (err error)
	inspectFuncRunReadCommitted   func(ctx context.Context, f func(ctxTX context.Context) error)
	afterRunReadCommittedCounter  uint64
	beforeRunReadCommittedCounter uint64
	RunReadCommittedMock          mTransactorMockRunReadCommitted

	funcRunRepeatableRead          func(ctx context.Context, f func(ctxTX context.Context) error) (err error)
	inspectFuncRunRepeatableRead   func(ctx context.Context, f func(ctxTX context.Context) error)
	afterRunRepeatableReadCounter  uint64
	beforeRunRepeatableReadCounter uint64
	RunRepeatableReadMock          mTransactorMockRunRepeatableRead

	funcRunSerializable          func(ctx context.Context, f func(ctxTX context.Context) error) (err error)
	inspectFuncRunSerializable   func(ctx context.Context, f func(ctxTX context.Context) error)
	afterRunSerializableCounter  uint64
	beforeRunSerializableCounter uint64
	RunSerializableMock          mTransactorMockRunSerializable
}

// NewTransactorMock returns a mock for domain.Transactor
func NewTransactorMock(t minimock.Tester) *TransactorMock {
	m := &TransactorMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.RunReadCommittedMock = mTransactorMockRunReadCommitted{mock: m}
	m.RunReadCommittedMock.callArgs = []*TransactorMockRunReadCommittedParams{}

	m.RunRepeatableReadMock = mTransactorMockRunRepeatableRead{mock: m}
	m.RunRepeatableReadMock.callArgs = []*TransactorMockRunRepeatableReadParams{}

	m.RunSerializableMock = mTransactorMockRunSerializable{mock: m}
	m.RunSerializableMock.callArgs = []*TransactorMockRunSerializableParams{}

	return m
}

type mTransactorMockRunReadCommitted struct {
	mock               *TransactorMock
	defaultExpectation *TransactorMockRunReadCommittedExpectation
	expectations       []*TransactorMockRunReadCommittedExpectation

	callArgs []*TransactorMockRunReadCommittedParams
	mutex    sync.RWMutex
}

// TransactorMockRunReadCommittedExpectation specifies expectation struct of the Transactor.RunReadCommitted
type TransactorMockRunReadCommittedExpectation struct {
	mock    *TransactorMock
	params  *TransactorMockRunReadCommittedParams
	results *TransactorMockRunReadCommittedResults
	Counter uint64
}

// TransactorMockRunReadCommittedParams contains parameters of the Transactor.RunReadCommitted
type TransactorMockRunReadCommittedParams struct {
	ctx context.Context
	f   func(ctxTX context.Context) error
}

// TransactorMockRunReadCommittedResults contains results of the Transactor.RunReadCommitted
type TransactorMockRunReadCommittedResults struct {
	err error
}

// Expect sets up expected params for Transactor.RunReadCommitted
func (mmRunReadCommitted *mTransactorMockRunReadCommitted) Expect(ctx context.Context, f func(ctxTX context.Context) error) *mTransactorMockRunReadCommitted {
	if mmRunReadCommitted.mock.funcRunReadCommitted != nil {
		mmRunReadCommitted.mock.t.Fatalf("TransactorMock.RunReadCommitted mock is already set by Set")
	}

	if mmRunReadCommitted.defaultExpectation == nil {
		mmRunReadCommitted.defaultExpectation = &TransactorMockRunReadCommittedExpectation{}
	}

	mmRunReadCommitted.defaultExpectation.params = &TransactorMockRunReadCommittedParams{ctx, f}
	for _, e := range mmRunReadCommitted.expectations {
		if minimock.Equal(e.params, mmRunReadCommitted.defaultExpectation.params) {
			mmRunReadCommitted.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmRunReadCommitted.defaultExpectation.params)
		}
	}

	return mmRunReadCommitted
}

// Inspect accepts an inspector function that has same arguments as the Transactor.RunReadCommitted
func (mmRunReadCommitted *mTransactorMockRunReadCommitted) Inspect(f func(ctx context.Context, f func(ctxTX context.Context) error)) *mTransactorMockRunReadCommitted {
	if mmRunReadCommitted.mock.inspectFuncRunReadCommitted != nil {
		mmRunReadCommitted.mock.t.Fatalf("Inspect function is already set for TransactorMock.RunReadCommitted")
	}

	mmRunReadCommitted.mock.inspectFuncRunReadCommitted = f

	return mmRunReadCommitted
}

// Return sets up results that will be returned by Transactor.RunReadCommitted
func (mmRunReadCommitted *mTransactorMockRunReadCommitted) Return(err error) *TransactorMock {
	if mmRunReadCommitted.mock.funcRunReadCommitted != nil {
		mmRunReadCommitted.mock.t.Fatalf("TransactorMock.RunReadCommitted mock is already set by Set")
	}

	if mmRunReadCommitted.defaultExpectation == nil {
		mmRunReadCommitted.defaultExpectation = &TransactorMockRunReadCommittedExpectation{mock: mmRunReadCommitted.mock}
	}
	mmRunReadCommitted.defaultExpectation.results = &TransactorMockRunReadCommittedResults{err}
	return mmRunReadCommitted.mock
}

// Set uses given function f to mock the Transactor.RunReadCommitted method
func (mmRunReadCommitted *mTransactorMockRunReadCommitted) Set(f func(ctx context.Context, f func(ctxTX context.Context) error) (err error)) *TransactorMock {
	if mmRunReadCommitted.defaultExpectation != nil {
		mmRunReadCommitted.mock.t.Fatalf("Default expectation is already set for the Transactor.RunReadCommitted method")
	}

	if len(mmRunReadCommitted.expectations) > 0 {
		mmRunReadCommitted.mock.t.Fatalf("Some expectations are already set for the Transactor.RunReadCommitted method")
	}

	mmRunReadCommitted.mock.funcRunReadCommitted = f
	return mmRunReadCommitted.mock
}

// When sets expectation for the Transactor.RunReadCommitted which will trigger the result defined by the following
// Then helper
func (mmRunReadCommitted *mTransactorMockRunReadCommitted) When(ctx context.Context, f func(ctxTX context.Context) error) *TransactorMockRunReadCommittedExpectation {
	if mmRunReadCommitted.mock.funcRunReadCommitted != nil {
		mmRunReadCommitted.mock.t.Fatalf("TransactorMock.RunReadCommitted mock is already set by Set")
	}

	expectation := &TransactorMockRunReadCommittedExpectation{
		mock:   mmRunReadCommitted.mock,
		params: &TransactorMockRunReadCommittedParams{ctx, f},
	}
	mmRunReadCommitted.expectations = append(mmRunReadCommitted.expectations, expectation)
	return expectation
}

// Then sets up Transactor.RunReadCommitted return parameters for the expectation previously defined by the When method
func (e *TransactorMockRunReadCommittedExpectation) Then(err error) *TransactorMock {
	e.results = &TransactorMockRunReadCommittedResults{err}
	return e.mock
}

// RunReadCommitted implements domain.Transactor
func (mmRunReadCommitted *TransactorMock) RunReadCommitted(ctx context.Context, f func(ctxTX context.Context) error) (err error) {
	mm_atomic.AddUint64(&mmRunReadCommitted.beforeRunReadCommittedCounter, 1)
	defer mm_atomic.AddUint64(&mmRunReadCommitted.afterRunReadCommittedCounter, 1)

	if mmRunReadCommitted.inspectFuncRunReadCommitted != nil {
		mmRunReadCommitted.inspectFuncRunReadCommitted(ctx, f)
	}

	mm_params := &TransactorMockRunReadCommittedParams{ctx, f}

	// Record call args
	mmRunReadCommitted.RunReadCommittedMock.mutex.Lock()
	mmRunReadCommitted.RunReadCommittedMock.callArgs = append(mmRunReadCommitted.RunReadCommittedMock.callArgs, mm_params)
	mmRunReadCommitted.RunReadCommittedMock.mutex.Unlock()

	for _, e := range mmRunReadCommitted.RunReadCommittedMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if mmRunReadCommitted.RunReadCommittedMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmRunReadCommitted.RunReadCommittedMock.defaultExpectation.Counter, 1)
		mm_want := mmRunReadCommitted.RunReadCommittedMock.defaultExpectation.params
		mm_got := TransactorMockRunReadCommittedParams{ctx, f}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmRunReadCommitted.t.Errorf("TransactorMock.RunReadCommitted got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmRunReadCommitted.RunReadCommittedMock.defaultExpectation.results
		if mm_results == nil {
			mmRunReadCommitted.t.Fatal("No results are set for the TransactorMock.RunReadCommitted")
		}
		return (*mm_results).err
	}
	if mmRunReadCommitted.funcRunReadCommitted != nil {
		return mmRunReadCommitted.funcRunReadCommitted(ctx, f)
	}
	mmRunReadCommitted.t.Fatalf("Unexpected call to TransactorMock.RunReadCommitted. %v %v", ctx, f)
	return
}

// RunReadCommittedAfterCounter returns a count of finished TransactorMock.RunReadCommitted invocations
func (mmRunReadCommitted *TransactorMock) RunReadCommittedAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmRunReadCommitted.afterRunReadCommittedCounter)
}

// RunReadCommittedBeforeCounter returns a count of TransactorMock.RunReadCommitted invocations
func (mmRunReadCommitted *TransactorMock) RunReadCommittedBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmRunReadCommitted.beforeRunReadCommittedCounter)
}

// Calls returns a list of arguments used in each call to TransactorMock.RunReadCommitted.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmRunReadCommitted *mTransactorMockRunReadCommitted) Calls() []*TransactorMockRunReadCommittedParams {
	mmRunReadCommitted.mutex.RLock()

	argCopy := make([]*TransactorMockRunReadCommittedParams, len(mmRunReadCommitted.callArgs))
	copy(argCopy, mmRunReadCommitted.callArgs)

	mmRunReadCommitted.mutex.RUnlock()

	return argCopy
}

// MinimockRunReadCommittedDone returns true if the count of the RunReadCommitted invocations corresponds
// the number of defined expectations
func (m *TransactorMock) MinimockRunReadCommittedDone() bool {
	for _, e := range m.RunReadCommittedMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.RunReadCommittedMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterRunReadCommittedCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcRunReadCommitted != nil && mm_atomic.LoadUint64(&m.afterRunReadCommittedCounter) < 1 {
		return false
	}
	return true
}

// MinimockRunReadCommittedInspect logs each unmet expectation
func (m *TransactorMock) MinimockRunReadCommittedInspect() {
	for _, e := range m.RunReadCommittedMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to TransactorMock.RunReadCommitted with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.RunReadCommittedMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterRunReadCommittedCounter) < 1 {
		if m.RunReadCommittedMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to TransactorMock.RunReadCommitted")
		} else {
			m.t.Errorf("Expected call to TransactorMock.RunReadCommitted with params: %#v", *m.RunReadCommittedMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcRunReadCommitted != nil && mm_atomic.LoadUint64(&m.afterRunReadCommittedCounter) < 1 {
		m.t.Error("Expected call to TransactorMock.RunReadCommitted")
	}
}

type mTransactorMockRunRepeatableRead struct {
	mock               *TransactorMock
	defaultExpectation *TransactorMockRunRepeatableReadExpectation
	expectations       []*TransactorMockRunRepeatableReadExpectation

	callArgs []*TransactorMockRunRepeatableReadParams
	mutex    sync.RWMutex
}

// TransactorMockRunRepeatableReadExpectation specifies expectation struct of the Transactor.RunRepeatableRead
type TransactorMockRunRepeatableReadExpectation struct {
	mock    *TransactorMock
	params  *TransactorMockRunRepeatableReadParams
	results *TransactorMockRunRepeatableReadResults
	Counter uint64
}

// TransactorMockRunRepeatableReadParams contains parameters of the Transactor.RunRepeatableRead
type TransactorMockRunRepeatableReadParams struct {
	ctx context.Context
	f   func(ctxTX context.Context) error
}

// TransactorMockRunRepeatableReadResults contains results of the Transactor.RunRepeatableRead
type TransactorMockRunRepeatableReadResults struct {
	err error
}

// Expect sets up expected params for Transactor.RunRepeatableRead
func (mmRunRepeatableRead *mTransactorMockRunRepeatableRead) Expect(ctx context.Context, f func(ctxTX context.Context) error) *mTransactorMockRunRepeatableRead {
	if mmRunRepeatableRead.mock.funcRunRepeatableRead != nil {
		mmRunRepeatableRead.mock.t.Fatalf("TransactorMock.RunRepeatableRead mock is already set by Set")
	}

	if mmRunRepeatableRead.defaultExpectation == nil {
		mmRunRepeatableRead.defaultExpectation = &TransactorMockRunRepeatableReadExpectation{}
	}

	mmRunRepeatableRead.defaultExpectation.params = &TransactorMockRunRepeatableReadParams{ctx, f}
	for _, e := range mmRunRepeatableRead.expectations {
		if minimock.Equal(e.params, mmRunRepeatableRead.defaultExpectation.params) {
			mmRunRepeatableRead.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmRunRepeatableRead.defaultExpectation.params)
		}
	}

	return mmRunRepeatableRead
}

// Inspect accepts an inspector function that has same arguments as the Transactor.RunRepeatableRead
func (mmRunRepeatableRead *mTransactorMockRunRepeatableRead) Inspect(f func(ctx context.Context, f func(ctxTX context.Context) error)) *mTransactorMockRunRepeatableRead {
	if mmRunRepeatableRead.mock.inspectFuncRunRepeatableRead != nil {
		mmRunRepeatableRead.mock.t.Fatalf("Inspect function is already set for TransactorMock.RunRepeatableRead")
	}

	mmRunRepeatableRead.mock.inspectFuncRunRepeatableRead = f

	return mmRunRepeatableRead
}

// Return sets up results that will be returned by Transactor.RunRepeatableRead
func (mmRunRepeatableRead *mTransactorMockRunRepeatableRead) Return(err error) *TransactorMock {
	if mmRunRepeatableRead.mock.funcRunRepeatableRead != nil {
		mmRunRepeatableRead.mock.t.Fatalf("TransactorMock.RunRepeatableRead mock is already set by Set")
	}

	if mmRunRepeatableRead.defaultExpectation == nil {
		mmRunRepeatableRead.defaultExpectation = &TransactorMockRunRepeatableReadExpectation{mock: mmRunRepeatableRead.mock}
	}
	mmRunRepeatableRead.defaultExpectation.results = &TransactorMockRunRepeatableReadResults{err}
	return mmRunRepeatableRead.mock
}

// Set uses given function f to mock the Transactor.RunRepeatableRead method
func (mmRunRepeatableRead *mTransactorMockRunRepeatableRead) Set(f func(ctx context.Context, f func(ctxTX context.Context) error) (err error)) *TransactorMock {
	if mmRunRepeatableRead.defaultExpectation != nil {
		mmRunRepeatableRead.mock.t.Fatalf("Default expectation is already set for the Transactor.RunRepeatableRead method")
	}

	if len(mmRunRepeatableRead.expectations) > 0 {
		mmRunRepeatableRead.mock.t.Fatalf("Some expectations are already set for the Transactor.RunRepeatableRead method")
	}

	mmRunRepeatableRead.mock.funcRunRepeatableRead = f
	return mmRunRepeatableRead.mock
}

// When sets expectation for the Transactor.RunRepeatableRead which will trigger the result defined by the following
// Then helper
func (mmRunRepeatableRead *mTransactorMockRunRepeatableRead) When(ctx context.Context, f func(ctxTX context.Context) error) *TransactorMockRunRepeatableReadExpectation {
	if mmRunRepeatableRead.mock.funcRunRepeatableRead != nil {
		mmRunRepeatableRead.mock.t.Fatalf("TransactorMock.RunRepeatableRead mock is already set by Set")
	}

	expectation := &TransactorMockRunRepeatableReadExpectation{
		mock:   mmRunRepeatableRead.mock,
		params: &TransactorMockRunRepeatableReadParams{ctx, f},
	}
	mmRunRepeatableRead.expectations = append(mmRunRepeatableRead.expectations, expectation)
	return expectation
}

// Then sets up Transactor.RunRepeatableRead return parameters for the expectation previously defined by the When method
func (e *TransactorMockRunRepeatableReadExpectation) Then(err error) *TransactorMock {
	e.results = &TransactorMockRunRepeatableReadResults{err}
	return e.mock
}

// RunRepeatableRead implements domain.Transactor
func (mmRunRepeatableRead *TransactorMock) RunRepeatableRead(ctx context.Context, f func(ctxTX context.Context) error) (err error) {
	mm_atomic.AddUint64(&mmRunRepeatableRead.beforeRunRepeatableReadCounter, 1)
	defer mm_atomic.AddUint64(&mmRunRepeatableRead.afterRunRepeatableReadCounter, 1)

	if mmRunRepeatableRead.inspectFuncRunRepeatableRead != nil {
		mmRunRepeatableRead.inspectFuncRunRepeatableRead(ctx, f)
	}

	mm_params := &TransactorMockRunRepeatableReadParams{ctx, f}

	// Record call args
	mmRunRepeatableRead.RunRepeatableReadMock.mutex.Lock()
	mmRunRepeatableRead.RunRepeatableReadMock.callArgs = append(mmRunRepeatableRead.RunRepeatableReadMock.callArgs, mm_params)
	mmRunRepeatableRead.RunRepeatableReadMock.mutex.Unlock()

	for _, e := range mmRunRepeatableRead.RunRepeatableReadMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if mmRunRepeatableRead.RunRepeatableReadMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmRunRepeatableRead.RunRepeatableReadMock.defaultExpectation.Counter, 1)
		mm_want := mmRunRepeatableRead.RunRepeatableReadMock.defaultExpectation.params
		mm_got := TransactorMockRunRepeatableReadParams{ctx, f}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmRunRepeatableRead.t.Errorf("TransactorMock.RunRepeatableRead got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmRunRepeatableRead.RunRepeatableReadMock.defaultExpectation.results
		if mm_results == nil {
			mmRunRepeatableRead.t.Fatal("No results are set for the TransactorMock.RunRepeatableRead")
		}
		return (*mm_results).err
	}
	if mmRunRepeatableRead.funcRunRepeatableRead != nil {
		return mmRunRepeatableRead.funcRunRepeatableRead(ctx, f)
	}
	mmRunRepeatableRead.t.Fatalf("Unexpected call to TransactorMock.RunRepeatableRead. %v %v", ctx, f)
	return
}

// RunRepeatableReadAfterCounter returns a count of finished TransactorMock.RunRepeatableRead invocations
func (mmRunRepeatableRead *TransactorMock) RunRepeatableReadAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmRunRepeatableRead.afterRunRepeatableReadCounter)
}

// RunRepeatableReadBeforeCounter returns a count of TransactorMock.RunRepeatableRead invocations
func (mmRunRepeatableRead *TransactorMock) RunRepeatableReadBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmRunRepeatableRead.beforeRunRepeatableReadCounter)
}

// Calls returns a list of arguments used in each call to TransactorMock.RunRepeatableRead.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmRunRepeatableRead *mTransactorMockRunRepeatableRead) Calls() []*TransactorMockRunRepeatableReadParams {
	mmRunRepeatableRead.mutex.RLock()

	argCopy := make([]*TransactorMockRunRepeatableReadParams, len(mmRunRepeatableRead.callArgs))
	copy(argCopy, mmRunRepeatableRead.callArgs)

	mmRunRepeatableRead.mutex.RUnlock()

	return argCopy
}

// MinimockRunRepeatableReadDone returns true if the count of the RunRepeatableRead invocations corresponds
// the number of defined expectations
func (m *TransactorMock) MinimockRunRepeatableReadDone() bool {
	for _, e := range m.RunRepeatableReadMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.RunRepeatableReadMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterRunRepeatableReadCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcRunRepeatableRead != nil && mm_atomic.LoadUint64(&m.afterRunRepeatableReadCounter) < 1 {
		return false
	}
	return true
}

// MinimockRunRepeatableReadInspect logs each unmet expectation
func (m *TransactorMock) MinimockRunRepeatableReadInspect() {
	for _, e := range m.RunRepeatableReadMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to TransactorMock.RunRepeatableRead with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.RunRepeatableReadMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterRunRepeatableReadCounter) < 1 {
		if m.RunRepeatableReadMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to TransactorMock.RunRepeatableRead")
		} else {
			m.t.Errorf("Expected call to TransactorMock.RunRepeatableRead with params: %#v", *m.RunRepeatableReadMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcRunRepeatableRead != nil && mm_atomic.LoadUint64(&m.afterRunRepeatableReadCounter) < 1 {
		m.t.Error("Expected call to TransactorMock.RunRepeatableRead")
	}
}

type mTransactorMockRunSerializable struct {
	mock               *TransactorMock
	defaultExpectation *TransactorMockRunSerializableExpectation
	expectations       []*TransactorMockRunSerializableExpectation

	callArgs []*TransactorMockRunSerializableParams
	mutex    sync.RWMutex
}

// TransactorMockRunSerializableExpectation specifies expectation struct of the Transactor.RunSerializable
type TransactorMockRunSerializableExpectation struct {
	mock    *TransactorMock
	params  *TransactorMockRunSerializableParams
	results *TransactorMockRunSerializableResults
	Counter uint64
}

// TransactorMockRunSerializableParams contains parameters of the Transactor.RunSerializable
type TransactorMockRunSerializableParams struct {
	ctx context.Context
	f   func(ctxTX context.Context) error
}

// TransactorMockRunSerializableResults contains results of the Transactor.RunSerializable
type TransactorMockRunSerializableResults struct {
	err error
}

// Expect sets up expected params for Transactor.RunSerializable
func (mmRunSerializable *mTransactorMockRunSerializable) Expect(ctx context.Context, f func(ctxTX context.Context) error) *mTransactorMockRunSerializable {
	if mmRunSerializable.mock.funcRunSerializable != nil {
		mmRunSerializable.mock.t.Fatalf("TransactorMock.RunSerializable mock is already set by Set")
	}

	if mmRunSerializable.defaultExpectation == nil {
		mmRunSerializable.defaultExpectation = &TransactorMockRunSerializableExpectation{}
	}

	mmRunSerializable.defaultExpectation.params = &TransactorMockRunSerializableParams{ctx, f}
	for _, e := range mmRunSerializable.expectations {
		if minimock.Equal(e.params, mmRunSerializable.defaultExpectation.params) {
			mmRunSerializable.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmRunSerializable.defaultExpectation.params)
		}
	}

	return mmRunSerializable
}

// Inspect accepts an inspector function that has same arguments as the Transactor.RunSerializable
func (mmRunSerializable *mTransactorMockRunSerializable) Inspect(f func(ctx context.Context, f func(ctxTX context.Context) error)) *mTransactorMockRunSerializable {
	if mmRunSerializable.mock.inspectFuncRunSerializable != nil {
		mmRunSerializable.mock.t.Fatalf("Inspect function is already set for TransactorMock.RunSerializable")
	}

	mmRunSerializable.mock.inspectFuncRunSerializable = f

	return mmRunSerializable
}

// Return sets up results that will be returned by Transactor.RunSerializable
func (mmRunSerializable *mTransactorMockRunSerializable) Return(err error) *TransactorMock {
	if mmRunSerializable.mock.funcRunSerializable != nil {
		mmRunSerializable.mock.t.Fatalf("TransactorMock.RunSerializable mock is already set by Set")
	}

	if mmRunSerializable.defaultExpectation == nil {
		mmRunSerializable.defaultExpectation = &TransactorMockRunSerializableExpectation{mock: mmRunSerializable.mock}
	}
	mmRunSerializable.defaultExpectation.results = &TransactorMockRunSerializableResults{err}
	return mmRunSerializable.mock
}

// Set uses given function f to mock the Transactor.RunSerializable method
func (mmRunSerializable *mTransactorMockRunSerializable) Set(f func(ctx context.Context, f func(ctxTX context.Context) error) (err error)) *TransactorMock {
	if mmRunSerializable.defaultExpectation != nil {
		mmRunSerializable.mock.t.Fatalf("Default expectation is already set for the Transactor.RunSerializable method")
	}

	if len(mmRunSerializable.expectations) > 0 {
		mmRunSerializable.mock.t.Fatalf("Some expectations are already set for the Transactor.RunSerializable method")
	}

	mmRunSerializable.mock.funcRunSerializable = f
	return mmRunSerializable.mock
}

// When sets expectation for the Transactor.RunSerializable which will trigger the result defined by the following
// Then helper
func (mmRunSerializable *mTransactorMockRunSerializable) When(ctx context.Context, f func(ctxTX context.Context) error) *TransactorMockRunSerializableExpectation {
	if mmRunSerializable.mock.funcRunSerializable != nil {
		mmRunSerializable.mock.t.Fatalf("TransactorMock.RunSerializable mock is already set by Set")
	}

	expectation := &TransactorMockRunSerializableExpectation{
		mock:   mmRunSerializable.mock,
		params: &TransactorMockRunSerializableParams{ctx, f},
	}
	mmRunSerializable.expectations = append(mmRunSerializable.expectations, expectation)
	return expectation
}

// Then sets up Transactor.RunSerializable return parameters for the expectation previously defined by the When method
func (e *TransactorMockRunSerializableExpectation) Then(err error) *TransactorMock {
	e.results = &TransactorMockRunSerializableResults{err}
	return e.mock
}

// RunSerializable implements domain.Transactor
func (mmRunSerializable *TransactorMock) RunSerializable(ctx context.Context, f func(ctxTX context.Context) error) (err error) {
	mm_atomic.AddUint64(&mmRunSerializable.beforeRunSerializableCounter, 1)
	defer mm_atomic.AddUint64(&mmRunSerializable.afterRunSerializableCounter, 1)

	if mmRunSerializable.inspectFuncRunSerializable != nil {
		mmRunSerializable.inspectFuncRunSerializable(ctx, f)
	}

	mm_params := &TransactorMockRunSerializableParams{ctx, f}

	// Record call args
	mmRunSerializable.RunSerializableMock.mutex.Lock()
	mmRunSerializable.RunSerializableMock.callArgs = append(mmRunSerializable.RunSerializableMock.callArgs, mm_params)
	mmRunSerializable.RunSerializableMock.mutex.Unlock()

	for _, e := range mmRunSerializable.RunSerializableMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if mmRunSerializable.RunSerializableMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmRunSerializable.RunSerializableMock.defaultExpectation.Counter, 1)
		mm_want := mmRunSerializable.RunSerializableMock.defaultExpectation.params
		mm_got := TransactorMockRunSerializableParams{ctx, f}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmRunSerializable.t.Errorf("TransactorMock.RunSerializable got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmRunSerializable.RunSerializableMock.defaultExpectation.results
		if mm_results == nil {
			mmRunSerializable.t.Fatal("No results are set for the TransactorMock.RunSerializable")
		}
		return (*mm_results).err
	}
	if mmRunSerializable.funcRunSerializable != nil {
		return mmRunSerializable.funcRunSerializable(ctx, f)
	}
	mmRunSerializable.t.Fatalf("Unexpected call to TransactorMock.RunSerializable. %v %v", ctx, f)
	return
}

// RunSerializableAfterCounter returns a count of finished TransactorMock.RunSerializable invocations
func (mmRunSerializable *TransactorMock) RunSerializableAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmRunSerializable.afterRunSerializableCounter)
}

// RunSerializableBeforeCounter returns a count of TransactorMock.RunSerializable invocations
func (mmRunSerializable *TransactorMock) RunSerializableBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmRunSerializable.beforeRunSerializableCounter)
}

// Calls returns a list of arguments used in each call to TransactorMock.RunSerializable.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmRunSerializable *mTransactorMockRunSerializable) Calls() []*TransactorMockRunSerializableParams {
	mmRunSerializable.mutex.RLock()

	argCopy := make([]*TransactorMockRunSerializableParams, len(mmRunSerializable.callArgs))
	copy(argCopy, mmRunSerializable.callArgs)

	mmRunSerializable.mutex.RUnlock()

	return argCopy
}

// MinimockRunSerializableDone returns true if the count of the RunSerializable invocations corresponds
// the number of defined expectations
func (m *TransactorMock) MinimockRunSerializableDone() bool {
	for _, e := range m.RunSerializableMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.RunSerializableMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterRunSerializableCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcRunSerializable != nil && mm_atomic.LoadUint64(&m.afterRunSerializableCounter) < 1 {
		return false
	}
	return true
}

// MinimockRunSerializableInspect logs each unmet expectation
func (m *TransactorMock) MinimockRunSerializableInspect() {
	for _, e := range m.RunSerializableMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to TransactorMock.RunSerializable with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.RunSerializableMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterRunSerializableCounter) < 1 {
		if m.RunSerializableMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to TransactorMock.RunSerializable")
		} else {
			m.t.Errorf("Expected call to TransactorMock.RunSerializable with params: %#v", *m.RunSerializableMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcRunSerializable != nil && mm_atomic.LoadUint64(&m.afterRunSerializableCounter) < 1 {
		m.t.Error("Expected call to TransactorMock.RunSerializable")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *TransactorMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockRunReadCommittedInspect()

		m.MinimockRunRepeatableReadInspect()

		m.MinimockRunSerializableInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *TransactorMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *TransactorMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockRunReadCommittedDone() &&
		m.MinimockRunRepeatableReadDone() &&
		m.MinimockRunSerializableDone()
}
