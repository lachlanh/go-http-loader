package main

import "fmt"

var V int

func F() { fmt.Printf("hello, number %d\n", V) }

func NewRequestBody() string {
	return "{'test':'hello'}"
}
