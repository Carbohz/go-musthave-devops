// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package storagemock is a generated GoMock package.
package storagemock

import (
	reflect "reflect"

	model "github.com/Carbohz/go-musthave-devops/model"
	gomock "github.com/golang/mock/gomock"
)

// MockMetricsStorager is a mock of MetricsStorager interface.
type MockMetricsStorager struct {
	ctrl     *gomock.Controller
	recorder *MockMetricsStoragerMockRecorder
}

// MockMetricsStoragerMockRecorder is the mock recorder for MockMetricsStorager.
type MockMetricsStoragerMockRecorder struct {
	mock *MockMetricsStorager
}

// NewMockMetricsStorager creates a new mock instance.
func NewMockMetricsStorager(ctrl *gomock.Controller) *MockMetricsStorager {
	mock := &MockMetricsStorager{ctrl: ctrl}
	mock.recorder = &MockMetricsStoragerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricsStorager) EXPECT() *MockMetricsStoragerMockRecorder {
	return m.recorder
}

// Dump mocks base method.
func (m *MockMetricsStorager) Dump() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Dump")
}

// Dump indicates an expected call of Dump.
func (mr *MockMetricsStoragerMockRecorder) Dump() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dump", reflect.TypeOf((*MockMetricsStorager)(nil).Dump))
}

// GetMetric mocks base method.
func (m *MockMetricsStorager) GetMetric(name string) (model.Metric, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetric", name)
	ret0, _ := ret[0].(model.Metric)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetMetric indicates an expected call of GetMetric.
func (mr *MockMetricsStoragerMockRecorder) GetMetric(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetric", reflect.TypeOf((*MockMetricsStorager)(nil).GetMetric), name)
}

// LoadOnStart mocks base method.
func (m *MockMetricsStorager) LoadOnStart() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "LoadOnStart")
}

// LoadOnStart indicates an expected call of LoadOnStart.
func (mr *MockMetricsStoragerMockRecorder) LoadOnStart() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadOnStart", reflect.TypeOf((*MockMetricsStorager)(nil).LoadOnStart))
}

// SaveMetric mocks base method.
func (m_2 *MockMetricsStorager) SaveMetric(m model.Metric) {
	m_2.ctrl.T.Helper()
	m_2.ctrl.Call(m_2, "SaveMetric", m)
}

// SaveMetric indicates an expected call of SaveMetric.
func (mr *MockMetricsStoragerMockRecorder) SaveMetric(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveMetric", reflect.TypeOf((*MockMetricsStorager)(nil).SaveMetric), m)
}
