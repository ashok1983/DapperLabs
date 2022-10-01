package main

import (
	"testing"
)

func TestHandleSignup(t *testing.T) {
	c := &gin.Context{}
	HandleSignup(c)
}
