// Copyright 2017 VMware, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/vmware/vic/lib/tether-ng (interfaces: Collector,Reporter,Signaler,Releaser,Waiter,Interactor,Reaper,Plugin,PluginRegistrar)

package mock_tether_ng

import (
	context "context"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	tether_ng "github.com/vmware/vic/lib/tether-ng"
	types "github.com/vmware/vic/lib/tether-ng/types"
)

// Mock of Collector interface
type MockCollector struct {
	ctrl     *gomock.Controller
	recorder *_MockCollectorRecorder
}

// Recorder for MockCollector (not exported)
type _MockCollectorRecorder struct {
	mock *MockCollector
}

func NewMockCollector(ctrl *gomock.Controller) *MockCollector {
	mock := &MockCollector{ctrl: ctrl}
	mock.recorder = &_MockCollectorRecorder{mock}
	return mock
}

func (_m *MockCollector) EXPECT() *_MockCollectorRecorder {
	return _m.recorder
}

func (_m *MockCollector) Collect(_param0 context.Context) {
	_m.ctrl.Call(_m, "Collect", _param0)
}

func (_mr *_MockCollectorRecorder) Collect(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Collect", arg0)
}

// Mock of Reporter interface
type MockReporter struct {
	ctrl     *gomock.Controller
	recorder *_MockReporterRecorder
}

// Recorder for MockReporter (not exported)
type _MockReporterRecorder struct {
	mock *MockReporter
}

func NewMockReporter(ctrl *gomock.Controller) *MockReporter {
	mock := &MockReporter{ctrl: ctrl}
	mock.recorder = &_MockReporterRecorder{mock}
	return mock
}

func (_m *MockReporter) EXPECT() *_MockReporterRecorder {
	return _m.recorder
}

func (_m *MockReporter) Report(_param0 context.Context, _param1 chan<- error) {
	_m.ctrl.Call(_m, "Report", _param0, _param1)
}

func (_mr *_MockReporterRecorder) Report(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Report", arg0, arg1)
}

// Mock of Signaler interface
type MockSignaler struct {
	ctrl     *gomock.Controller
	recorder *_MockSignalerRecorder
}

// Recorder for MockSignaler (not exported)
type _MockSignalerRecorder struct {
	mock *MockSignaler
}

func NewMockSignaler(ctrl *gomock.Controller) *MockSignaler {
	mock := &MockSignaler{ctrl: ctrl}
	mock.recorder = &_MockSignalerRecorder{mock}
	return mock
}

func (_m *MockSignaler) EXPECT() *_MockSignalerRecorder {
	return _m.recorder
}

func (_m *MockSignaler) Kill(_param0 context.Context, _param1 string) error {
	ret := _m.ctrl.Call(_m, "Kill", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockSignalerRecorder) Kill(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Kill", arg0, arg1)
}

func (_m *MockSignaler) Running(_param0 context.Context, _param1 string) bool {
	ret := _m.ctrl.Call(_m, "Running", _param0, _param1)
	ret0, _ := ret[0].(bool)
	return ret0
}

func (_mr *_MockSignalerRecorder) Running(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Running", arg0, arg1)
}

// Mock of Releaser interface
type MockReleaser struct {
	ctrl     *gomock.Controller
	recorder *_MockReleaserRecorder
}

// Recorder for MockReleaser (not exported)
type _MockReleaserRecorder struct {
	mock *MockReleaser
}

func NewMockReleaser(ctrl *gomock.Controller) *MockReleaser {
	mock := &MockReleaser{ctrl: ctrl}
	mock.recorder = &_MockReleaserRecorder{mock}
	return mock
}

func (_m *MockReleaser) EXPECT() *_MockReleaserRecorder {
	return _m.recorder
}

func (_m *MockReleaser) Release(_param0 context.Context, _param1 chan<- chan struct{}) {
	_m.ctrl.Call(_m, "Release", _param0, _param1)
}

func (_mr *_MockReleaserRecorder) Release(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Release", arg0, arg1)
}

// Mock of Waiter interface
type MockWaiter struct {
	ctrl     *gomock.Controller
	recorder *_MockWaiterRecorder
}

// Recorder for MockWaiter (not exported)
type _MockWaiterRecorder struct {
	mock *MockWaiter
}

func NewMockWaiter(ctrl *gomock.Controller) *MockWaiter {
	mock := &MockWaiter{ctrl: ctrl}
	mock.recorder = &_MockWaiterRecorder{mock}
	return mock
}

func (_m *MockWaiter) EXPECT() *_MockWaiterRecorder {
	return _m.recorder
}

func (_m *MockWaiter) Wait(_param0 context.Context, _param1 <-chan chan struct{}) {
	_m.ctrl.Call(_m, "Wait", _param0, _param1)
}

func (_mr *_MockWaiterRecorder) Wait(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Wait", arg0, arg1)
}

// Mock of Interactor interface
type MockInteractor struct {
	ctrl     *gomock.Controller
	recorder *_MockInteractorRecorder
}

// Recorder for MockInteractor (not exported)
type _MockInteractorRecorder struct {
	mock *MockInteractor
}

func NewMockInteractor(ctrl *gomock.Controller) *MockInteractor {
	mock := &MockInteractor{ctrl: ctrl}
	mock.recorder = &_MockInteractorRecorder{mock}
	return mock
}

func (_m *MockInteractor) EXPECT() *_MockInteractorRecorder {
	return _m.recorder
}

func (_m *MockInteractor) Close(_param0 context.Context, _param1 <-chan *types.Session) <-chan struct{} {
	ret := _m.ctrl.Call(_m, "Close", _param0, _param1)
	ret0, _ := ret[0].(<-chan struct{})
	return ret0
}

func (_mr *_MockInteractorRecorder) Close(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Close", arg0, arg1)
}

func (_m *MockInteractor) NonInteract(_param0 context.Context, _param1 <-chan *types.Session) <-chan struct{} {
	ret := _m.ctrl.Call(_m, "NonInteract", _param0, _param1)
	ret0, _ := ret[0].(<-chan struct{})
	return ret0
}

func (_mr *_MockInteractorRecorder) NonInteract(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "NonInteract", arg0, arg1)
}

func (_m *MockInteractor) PseudoTerminal(_param0 context.Context, _param1 <-chan *types.Session) <-chan struct{} {
	ret := _m.ctrl.Call(_m, "PseudoTerminal", _param0, _param1)
	ret0, _ := ret[0].(<-chan struct{})
	return ret0
}

func (_mr *_MockInteractorRecorder) PseudoTerminal(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "PseudoTerminal", arg0, arg1)
}

// Mock of Reaper interface
type MockReaper struct {
	ctrl     *gomock.Controller
	recorder *_MockReaperRecorder
}

// Recorder for MockReaper (not exported)
type _MockReaperRecorder struct {
	mock *MockReaper
}

func NewMockReaper(ctrl *gomock.Controller) *MockReaper {
	mock := &MockReaper{ctrl: ctrl}
	mock.recorder = &_MockReaperRecorder{mock}
	return mock
}

func (_m *MockReaper) EXPECT() *_MockReaperRecorder {
	return _m.recorder
}

func (_m *MockReaper) Reap(_param0 context.Context) error {
	ret := _m.ctrl.Call(_m, "Reap", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockReaperRecorder) Reap(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Reap", arg0)
}

// Mock of Plugin interface
type MockPlugin struct {
	ctrl     *gomock.Controller
	recorder *_MockPluginRecorder
}

// Recorder for MockPlugin (not exported)
type _MockPluginRecorder struct {
	mock *MockPlugin
}

func NewMockPlugin(ctrl *gomock.Controller) *MockPlugin {
	mock := &MockPlugin{ctrl: ctrl}
	mock.recorder = &_MockPluginRecorder{mock}
	return mock
}

func (_m *MockPlugin) EXPECT() *_MockPluginRecorder {
	return _m.recorder
}

func (_m *MockPlugin) Configure(_param0 context.Context, _param1 *types.ExecutorConfig) error {
	ret := _m.ctrl.Call(_m, "Configure", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockPluginRecorder) Configure(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Configure", arg0, arg1)
}

func (_m *MockPlugin) Start(_param0 context.Context) error {
	ret := _m.ctrl.Call(_m, "Start", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockPluginRecorder) Start(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Start", arg0)
}

func (_m *MockPlugin) Stop(_param0 context.Context) error {
	ret := _m.ctrl.Call(_m, "Stop", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockPluginRecorder) Stop(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Stop", arg0)
}

func (_m *MockPlugin) UUID(_param0 context.Context) uuid.UUID {
	ret := _m.ctrl.Call(_m, "UUID", _param0)
	ret0, _ := ret[0].(uuid.UUID)
	return ret0
}

func (_mr *_MockPluginRecorder) UUID(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UUID", arg0)
}

// Mock of PluginRegistrar interface
type MockPluginRegistrar struct {
	ctrl     *gomock.Controller
	recorder *_MockPluginRegistrarRecorder
}

// Recorder for MockPluginRegistrar (not exported)
type _MockPluginRegistrarRecorder struct {
	mock *MockPluginRegistrar
}

func NewMockPluginRegistrar(ctrl *gomock.Controller) *MockPluginRegistrar {
	mock := &MockPluginRegistrar{ctrl: ctrl}
	mock.recorder = &_MockPluginRegistrarRecorder{mock}
	return mock
}

func (_m *MockPluginRegistrar) EXPECT() *_MockPluginRegistrarRecorder {
	return _m.recorder
}

func (_m *MockPluginRegistrar) Plugins(_param0 context.Context) []tether_ng.Plugin {
	ret := _m.ctrl.Call(_m, "Plugins", _param0)
	ret0, _ := ret[0].([]tether_ng.Plugin)
	return ret0
}

func (_mr *_MockPluginRegistrarRecorder) Plugins(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Plugins", arg0)
}

func (_m *MockPluginRegistrar) Register(_param0 context.Context, _param1 tether_ng.Plugin) error {
	ret := _m.ctrl.Call(_m, "Register", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockPluginRegistrarRecorder) Register(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Register", arg0, arg1)
}

func (_m *MockPluginRegistrar) Unregister(_param0 context.Context, _param1 tether_ng.Plugin) error {
	ret := _m.ctrl.Call(_m, "Unregister", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockPluginRegistrarRecorder) Unregister(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Unregister", arg0, arg1)
}
