// Code generated by mockery v1.0.0
package mocks

import mock "github.com/stretchr/testify/mock"
import model "github.com/getamis/eth-indexer/model"

// Store is an autogenerated mock type for the Store type
type Store struct {
	mock.Mock
}

// DeleteByBlockNumber provides a mock function with given fields: from, to
func (_m *Store) DeleteByBlockNumber(from int64, to int64) error {
	ret := _m.Called(from, to)

	var r0 error
	if rf, ok := ret.Get(0).(func(int64, int64) error); ok {
		r0 = rf(from, to)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindUncleByHash provides a mock function with given fields: hash
func (_m *Store) FindUncleByHash(hash []byte) (*model.UncleHeader, error) {
	ret := _m.Called(hash)

	var r0 *model.UncleHeader
	if rf, ok := ret.Get(0).(func([]byte) *model.UncleHeader); ok {
		r0 = rf(hash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.UncleHeader)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(hash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Insert provides a mock function with given fields: data
func (_m *Store) Insert(data *model.UncleHeader) error {
	ret := _m.Called(data)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.UncleHeader) error); ok {
		r0 = rf(data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}