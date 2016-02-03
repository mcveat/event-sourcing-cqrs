package store

import (
	. "github.com/go-check/check"
	"testing"
)

func Test(t *testing.T) {
	TestingT(t)
}

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestSaveEvent(c *C) {
	es := Empty()
	uuid := es.Save(1)
	c.Assert(uuid, NotNil)
	c.Assert(es, HasLen, 1)
	c.Assert(es[(*uuid)], Equals, 1)
}
