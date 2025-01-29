package frn_test

import (
	"github.com/Freemodel-Inc/frn"
	"github.com/segmentio/ksuid"
	"github.com/tj/assert"
	"testing"
)

type MockHandler struct {
	randomString func() string
}

func defaultRandomString() string {
	return ksuid.New().String()
}

func newMockHandler() *MockHandler {
	return &MockHandler{
		randomString: defaultRandomString,
	}
}

// WithRandomString is for testing to allow randomString to be overridden
func (m *MockHandler) WithRandomString(fn func() string) *MockHandler {
	m.randomString = fn
	return m
}

func (m *MockHandler) DoWork(projectID frn.ID) (contract frn.ID) {
	return projectID.Sub("contract", m.randomString())
}

func TestNewSequence(t *testing.T) {
	var (
		seq       = frn.NewSequence(1)
		handler   = newMockHandler().WithRandomString(seq.Next)
		projectID = frn.ID("dev:crm:project:1")
	)

	contractID := handler.DoWork(projectID)
	assert.EqualValues(t, "dev:crm:project:1:contract:2", contractID)
}
