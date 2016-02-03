package main

import (
  "testing"
  . "github.com/go-check/check"
)

func Test(t *testing.T) {
  TestingT(t)
}

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestHelloWorld(c *C) {
  c.Check(42, Equals, 42)
}
