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

func (s *MySuite) TestSaveEvents(c *C) {
	es := Empty()
	uuid := es.Save([]Event{GenericEvent{value: 42}, GenericEvent{value: -1}})
	c.Assert(uuid, NotNil)
	c.Assert(es, HasLen, 1)
	events := es[(*uuid)]
	c.Assert(events, HasLen, 2)
	c.Assert(events[0], Equals, GenericEvent{uuid: uuid, value: 42})
	c.Assert(events[1], Equals, GenericEvent{uuid: uuid, value: -1})
}
