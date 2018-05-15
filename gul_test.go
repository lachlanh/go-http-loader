package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type HttpClientMock struct {
}

func (c *HttpClientMock) Get(url string) (*http.Response, error) {
	return &http.Response{}, nil
}

func TestRunUrl(t *testing.T) {
	clientMock := HttpClientMock{}
	result := runUrl("http://127.0.0.1", 10, &clientMock)
	assert.NotNil(t, result)
	assert.Equal(t, 10, result.successCount)
}
