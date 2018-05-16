package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"plugin"
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

func TestLoadPlugin(t *testing.T) {
	cmd := exec.Command("go", "build", "-buildmode=plugin", "request_plg.go")
	defer os.Remove("request_plg.so")
	var out, errout bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errout
	err := cmd.Run()
	if err != nil {
		fmt.Printf("err : %q\n", errout.String())
		panic(err)
	}
	fmt.Printf("out : %q\n", out.String())

	p, err := plugin.Open("request_plg.so")
	if err != nil {
		panic(err)
	}
	v, err := p.Lookup("V")
	if err != nil {
		panic(err)
	}
	f, err := p.Lookup("F")
	if err != nil {
		panic(err)
	}
	*v.(*int) = 7
	f.(func())()
	assert.Fail(t, "fail")
}
