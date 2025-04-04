// Copyright (c) 2025 Sidero Labs, Inc.
//
// Use of this software is governed by the Business Source License
// included in the LICENSE file.
//

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/siderolabs/omni/internal/backend/runtime/omni/controllers/omni/etcdbackup (interfaces: Store)
//
// Generated by this command:
//
//	mockgen -destination=etcd_backup_mock_3_test.go -package omni_test -typed -copyright_file ../../../../../../hack/.license-header.go.txt github.com/siderolabs/omni/internal/backend/runtime/omni/controllers/omni/etcdbackup Store
//

// Package omni_test is a generated GoMock package.
package omni_test

import (
	context "context"
	io "io"
	iter "iter"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"

	etcdbackup "github.com/siderolabs/omni/internal/backend/runtime/omni/controllers/omni/etcdbackup"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
	isgomock struct{}
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// Download mocks base method.
func (m *MockStore) Download(ctx context.Context, encryptionKey []byte, clusterUUID, snapshotName string) (etcdbackup.BackupData, io.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Download", ctx, encryptionKey, clusterUUID, snapshotName)
	ret0, _ := ret[0].(etcdbackup.BackupData)
	ret1, _ := ret[1].(io.ReadCloser)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Download indicates an expected call of Download.
func (mr *MockStoreMockRecorder) Download(ctx, encryptionKey, clusterUUID, snapshotName any) *MockStoreDownloadCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Download", reflect.TypeOf((*MockStore)(nil).Download), ctx, encryptionKey, clusterUUID, snapshotName)
	return &MockStoreDownloadCall{Call: call}
}

// MockStoreDownloadCall wrap *gomock.Call
type MockStoreDownloadCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStoreDownloadCall) Return(arg0 etcdbackup.BackupData, arg1 io.ReadCloser, arg2 error) *MockStoreDownloadCall {
	c.Call = c.Call.Return(arg0, arg1, arg2)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStoreDownloadCall) Do(f func(context.Context, []byte, string, string) (etcdbackup.BackupData, io.ReadCloser, error)) *MockStoreDownloadCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStoreDownloadCall) DoAndReturn(f func(context.Context, []byte, string, string) (etcdbackup.BackupData, io.ReadCloser, error)) *MockStoreDownloadCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ListBackups mocks base method.
func (m *MockStore) ListBackups(ctx context.Context, clusterUUID string) (iter.Seq2[etcdbackup.Info, error], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListBackups", ctx, clusterUUID)
	ret0, _ := ret[0].(iter.Seq2[etcdbackup.Info, error])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListBackups indicates an expected call of ListBackups.
func (mr *MockStoreMockRecorder) ListBackups(ctx, clusterUUID any) *MockStoreListBackupsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListBackups", reflect.TypeOf((*MockStore)(nil).ListBackups), ctx, clusterUUID)
	return &MockStoreListBackupsCall{Call: call}
}

// MockStoreListBackupsCall wrap *gomock.Call
type MockStoreListBackupsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStoreListBackupsCall) Return(arg0 iter.Seq2[etcdbackup.Info, error], arg1 error) *MockStoreListBackupsCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStoreListBackupsCall) Do(f func(context.Context, string) (iter.Seq2[etcdbackup.Info, error], error)) *MockStoreListBackupsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStoreListBackupsCall) DoAndReturn(f func(context.Context, string) (iter.Seq2[etcdbackup.Info, error], error)) *MockStoreListBackupsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Upload mocks base method.
func (m *MockStore) Upload(ctx context.Context, descr etcdbackup.Description, r io.Reader) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Upload", ctx, descr, r)
	ret0, _ := ret[0].(error)
	return ret0
}

// Upload indicates an expected call of Upload.
func (mr *MockStoreMockRecorder) Upload(ctx, descr, r any) *MockStoreUploadCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upload", reflect.TypeOf((*MockStore)(nil).Upload), ctx, descr, r)
	return &MockStoreUploadCall{Call: call}
}

// MockStoreUploadCall wrap *gomock.Call
type MockStoreUploadCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStoreUploadCall) Return(arg0 error) *MockStoreUploadCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStoreUploadCall) Do(f func(context.Context, etcdbackup.Description, io.Reader) error) *MockStoreUploadCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStoreUploadCall) DoAndReturn(f func(context.Context, etcdbackup.Description, io.Reader) error) *MockStoreUploadCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
