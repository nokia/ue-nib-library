/*
   Copyright (c) 2020 Nokia.

   Licensed under the BSD 3-Clause Clear License.
   SPDX-License-Identifier: BSD-3-Clause-Clear
*/

package uenibreader_test

import (
	"errors"
	"github.com/nokia/ue-nib-library/pkg/uenibreader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockSdlBackend struct {
	mock.Mock
}

func (m *mockSdlBackend) Get(ns string, keys []string) (map[string]interface{}, error) {
	a := m.Called(ns, keys)
	if a.Get(0) == nil {
		return nil, a.Error(1)
	}
	return a.Get(0).(map[string]interface{}), a.Error(1)
}

func (m *mockSdlBackend) SubscribeChannel(ns string, cb func(string, ...string), channels ...string) error {
	a := m.Called(ns, cb, channels)
	return a.Error(0)
}

func (m *mockSdlBackend) Close() error {
	a := m.Called()
	return a.Error(0)
}

func setup() (*mockSdlBackend, *uenibreader.Reader) {
	uenibreader.SetDisableSdlCreationInConstructor(true)
	m := new(mockSdlBackend)
	i := uenibreader.NewReader()
	i.SetDbBackend(m)
	return m, i
}

func TestCanCreateReaderInstance(t *testing.T) {
	_, i := setup()
	assert.NotNil(t, i)
}

func TestCloseSuccess(t *testing.T) {
	m, i := setup()
	m.On("Close").Return(nil)
	err := i.Close()
	assert.Nil(t, err)
	m.AssertExpectations(t)
}

func TestCloseDbBackendFailure(t *testing.T) {
	m, i := setup()
	m.On("Close").Return(errors.New("Some DB Backend Error"))
	err := i.Close()
	uenibFailure, ok := err.(uenibreader.Error)
	assert.Equal(t, true, ok)
	assert.Equal(t, true, uenibFailure.Temporary())
	assert.Contains(t, uenibFailure.Error(), "database backend error: Some DB Backend Error")
	m.AssertExpectations(t)
}
