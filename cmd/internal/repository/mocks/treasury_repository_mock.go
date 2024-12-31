// Code generated by MockGen. DO NOT EDIT.
// Source: ./treasury_repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/pablorodrigo52/transaction-api/cmd/internal/model"
)

// MockTreasuryRepository is a mock of TreasuryRepository interface.
type MockTreasuryRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTreasuryRepositoryMockRecorder
}

// MockTreasuryRepositoryMockRecorder is the mock recorder for MockTreasuryRepository.
type MockTreasuryRepositoryMockRecorder struct {
	mock *MockTreasuryRepository
}

// NewMockTreasuryRepository creates a new mock instance.
func NewMockTreasuryRepository(ctrl *gomock.Controller) *MockTreasuryRepository {
	mock := &MockTreasuryRepository{ctrl: ctrl}
	mock.recorder = &MockTreasuryRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTreasuryRepository) EXPECT() *MockTreasuryRepositoryMockRecorder {
	return m.recorder
}

// GetExchangeRateByCountry mocks base method.
func (m *MockTreasuryRepository) GetExchangeRateByCountry(country string) (*model.TreasuryRatesExchange, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetExchangeRateByCountry", country)
	ret0, _ := ret[0].(*model.TreasuryRatesExchange)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetExchangeRateByCountry indicates an expected call of GetExchangeRateByCountry.
func (mr *MockTreasuryRepositoryMockRecorder) GetExchangeRateByCountry(country interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetExchangeRateByCountry", reflect.TypeOf((*MockTreasuryRepository)(nil).GetExchangeRateByCountry), country)
}
