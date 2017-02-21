// Copyright 2016 Mender Software AS
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

// generate with:
// mockery -name=DevAuthApp  -inpkg -print
package main

import log "github.com/mendersoftware/deviceauth/log"
import mock "github.com/stretchr/testify/mock"

// MockDevAuthApp is an autogenerated mock type for the DevAuthApp type
type MockDevAuthApp struct {
	mock.Mock
}

// AcceptDevice provides a mock function with given fields: dev_id
func (_m *MockDevAuthApp) AcceptDevice(dev_id string) error {
	ret := _m.Called(dev_id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(dev_id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetDevice provides a mock function with given fields: dev_id
func (_m *MockDevAuthApp) GetDevice(dev_id string) (*Device, error) {
	ret := _m.Called(dev_id)

	var r0 *Device
	if rf, ok := ret.Get(0).(func(string) *Device); ok {
		r0 = rf(dev_id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Device)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(dev_id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDeviceToken provides a mock function with given fields: dev_id
func (_m *MockDevAuthApp) GetDeviceToken(dev_id string) (*Token, error) {
	ret := _m.Called(dev_id)

	var r0 *Token
	if rf, ok := ret.Get(0).(func(string) *Token); ok {
		r0 = rf(dev_id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Token)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(dev_id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDevices provides a mock function with given fields: skip, limit
func (_m *MockDevAuthApp) GetDevices(skip uint, limit uint) ([]Device, error) {
	ret := _m.Called(skip, limit)

	var r0 []Device
	if rf, ok := ret.Get(0).(func(uint, uint) []Device); ok {
		r0 = rf(skip, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]Device)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint, uint) error); ok {
		r1 = rf(skip, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RejectDevice provides a mock function with given fields: dev_id
func (_m *MockDevAuthApp) RejectDevice(dev_id string) error {
	ret := _m.Called(dev_id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(dev_id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ResetDevice provides a mock function with given fields: dev_id
func (_m *MockDevAuthApp) ResetDevice(dev_id string) error {
	ret := _m.Called(dev_id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(dev_id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RevokeToken provides a mock function with given fields: token_id
func (_m *MockDevAuthApp) RevokeToken(token_id string) error {
	ret := _m.Called(token_id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(token_id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SubmitAuthRequest provides a mock function with given fields: r
func (_m *MockDevAuthApp) SubmitAuthRequest(r *AuthReq) (string, error) {
	ret := _m.Called(r)

	var r0 string
	if rf, ok := ret.Get(0).(func(*AuthReq) string); ok {
		r0 = rf(r)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*AuthReq) error); ok {
		r1 = rf(r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UseLog provides a mock function with given fields: l
func (_m *MockDevAuthApp) UseLog(l *log.Logger) {
	_m.Called(l)
}

// VerifyToken provides a mock function with given fields: token
func (_m *MockDevAuthApp) VerifyToken(token string) error {
	ret := _m.Called(token)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WithContext provides a mock function with given fields: c
func (_m *MockDevAuthApp) WithContext(c *RequestContext) DevAuthApp {
	ret := _m.Called(c)

	var r0 DevAuthApp
	if rf, ok := ret.Get(0).(func(*RequestContext) DevAuthApp); ok {
		r0 = rf(c)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(DevAuthApp)
		}
	}

	return r0
}

var _ DevAuthApp = (*MockDevAuthApp)(nil)
