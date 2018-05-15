package main

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type HttpClientMock struct {
}

func (c *HttpClientMock) Get(url string) (*http.Response, error) {
	return &http.Response{}, nil
}

func TestRunUrl(t *testing.T) {
	clientMock := HttpClientMock{}
	result := runUrl("http://127.0.0.1", 10, 10, &clientMock)
	assert.NotNil(t, result)
	assert.Equal(t, 10, result.successCount)
}

func TestCheckTiming(t *testing.T) {
	timeLimit, _ := time.ParseDuration("1s")
	start := time.Now()
	time.Sleep(timeLimit)
	elapsed := time.Since(start)
	fmt.Println("elapsed : ", elapsed)
	if time.Since(start) >= timeLimit {
		fmt.Println("here we are")
	}
	assert.Fail(t, "fail")

}
