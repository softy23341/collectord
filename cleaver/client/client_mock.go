package client

import (
	"sync"

	"git.softndit.com/collector/backend/cleaver"
)

type mockResizeRV struct {
	results []*cleaver.TransformResult
	err     error
}

type mockCopyRV struct {
	result *cleaver.CopyResult
	err    error
}

var mockErrGetNotFound = cleaver.ErrCategoryGetNotFound.New("mock client unregistered source")

// NewMockClient TBD
func NewMockClient() *MockClient {
	return &MockClient{
		resizeRVBySource: make(map[string]mockResizeRV),
		copyRVBySource:   make(map[string]mockCopyRV),
	}
}

// Check that MockClient implements Client interface
var _ Client = &MockClient{}

// MockClient TBD
type MockClient struct {
	sync.Mutex

	resizeRVBySource map[string]mockResizeRV
	copyRVBySource   map[string]mockCopyRV
}

// AddResizeRVBySource TBD
func (m *MockClient) AddResizeRVBySource(source string, results []*cleaver.TransformResult, err error) {
	m.Lock()
	defer m.Unlock()
	m.resizeRVBySource[source] = mockResizeRV{results: results, err: err}
}

// AddCopyRVBySource TBD
func (m *MockClient) AddCopyRVBySource(source string, result *cleaver.CopyResult, err error) {
	m.Lock()
	defer m.Unlock()
	m.copyRVBySource[source] = mockCopyRV{result: result, err: err}
}

// Resize TBD
func (m *MockClient) Resize(task *cleaver.ResizeTask) ([]*cleaver.TransformResult, error) {
	m.Lock()
	defer m.Unlock()
	rv, ok := m.resizeRVBySource[task.Source]
	if !ok {
		return nil, mockErrGetNotFound
	}
	return rv.results, rv.err
}

// Copy TBD
func (m *MockClient) Copy(task *cleaver.CopyTask) (*cleaver.CopyResult, error) {
	m.Lock()
	defer m.Unlock()
	rv, ok := m.copyRVBySource[task.Source]
	if !ok {
		return nil, mockErrGetNotFound
	}
	return rv.result, rv.err
}
