// Code generated by MockGen. DO NOT EDIT.
// Source: ./transaction_repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/pablorodrigo52/transaction-api/cmd/internal/model"
)

// MockTransactionRepository is a mock of TransactionRepository interface.
type MockTransactionRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionRepositoryMockRecorder
}

// MockTransactionRepositoryMockRecorder is the mock recorder for MockTransactionRepository.
type MockTransactionRepositoryMockRecorder struct {
	mock *MockTransactionRepository
}

// NewMockTransactionRepository creates a new mock instance.
func NewMockTransactionRepository(ctrl *gomock.Controller) *MockTransactionRepository {
	mock := &MockTransactionRepository{ctrl: ctrl}
	mock.recorder = &MockTransactionRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactionRepository) EXPECT() *MockTransactionRepositoryMockRecorder {
	return m.recorder
}

// GetTransaction mocks base method.
func (m *MockTransactionRepository) GetTransaction(transactionID int64) (*model.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransaction", transactionID)
	ret0, _ := ret[0].(*model.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransaction indicates an expected call of GetTransaction.
func (mr *MockTransactionRepositoryMockRecorder) GetTransaction(transactionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransaction", reflect.TypeOf((*MockTransactionRepository)(nil).GetTransaction), transactionID)
}

// LogicalDeleteTransaction mocks base method.
func (m *MockTransactionRepository) LogicalDeleteTransaction(transactionID int64) (*int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogicalDeleteTransaction", transactionID)
	ret0, _ := ret[0].(*int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LogicalDeleteTransaction indicates an expected call of LogicalDeleteTransaction.
func (mr *MockTransactionRepositoryMockRecorder) LogicalDeleteTransaction(transactionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogicalDeleteTransaction", reflect.TypeOf((*MockTransactionRepository)(nil).LogicalDeleteTransaction), transactionID)
}

// SaveTransaction mocks base method.
func (m *MockTransactionRepository) SaveTransaction(transaction *model.Transaction) (*model.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveTransaction", transaction)
	ret0, _ := ret[0].(*model.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveTransaction indicates an expected call of SaveTransaction.
func (mr *MockTransactionRepositoryMockRecorder) SaveTransaction(transaction interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveTransaction", reflect.TypeOf((*MockTransactionRepository)(nil).SaveTransaction), transaction)
}

// UpdateTransaction mocks base method.
func (m *MockTransactionRepository) UpdateTransaction(transactionID int64, transaction *model.Transaction) (*model.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTransaction", transactionID, transaction)
	ret0, _ := ret[0].(*model.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateTransaction indicates an expected call of UpdateTransaction.
func (mr *MockTransactionRepositoryMockRecorder) UpdateTransaction(transactionID, transaction interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTransaction", reflect.TypeOf((*MockTransactionRepository)(nil).UpdateTransaction), transactionID, transaction)
}
