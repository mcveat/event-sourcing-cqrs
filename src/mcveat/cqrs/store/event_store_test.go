package store

import (
	. "github.com/go-check/check"
	. "github.com/nu7hatch/gouuid"
	"testing"
)

func Test(t *testing.T) {
	TestingT(t)
}

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestSaveEventsInStore(c *C) {
	es := Empty()
	uuid := es.Save([]Event{GenericEvent{value: 42}, GenericEvent{value: -1}})
	c.Assert(uuid, NotNil)
	c.Assert(es.store, HasLen, 1)
	events := es.store[(*uuid)]
	c.Assert(events, HasLen, 2)
	c.Assert(events[0], Equals, GenericEvent{uuid: uuid, value: 42})
	c.Assert(events[1], Equals, GenericEvent{uuid: uuid, value: -1})
}

func (s *MySuite) TestSaveEventsInLog(c *C) {
	es := Empty()
	firstUuid := es.Save([]Event{GenericEvent{value: 1}, GenericEvent{value: 2}})
	secondUuid := es.Save([]Event{GenericEvent{value: 3}, GenericEvent{value: 4}})
	c.Assert(firstUuid, Not(Equals), secondUuid)
	log := es.log
	c.Assert(log, HasLen, 4)
	c.Assert(log[0], Equals, GenericEvent{uuid: firstUuid, value: 1})
	c.Assert(log[1], Equals, GenericEvent{uuid: firstUuid, value: 2})
	c.Assert(log[2], Equals, GenericEvent{uuid: secondUuid, value: 3})
	c.Assert(log[3], Equals, GenericEvent{uuid: secondUuid, value: 4})
}

func (s *MySuite) TestFindMissing(c *C) {
	es := Empty()
	randomUUID, _ := NewV4()
	history := es.Find(randomUUID)
	c.Assert(history.events, HasLen, 0)
	c.Assert(history.version, Equals, 0)
}

func (s *MySuite) TestFind(c *C) {
	es := Empty()
	uuid := es.Save([]Event{GenericEvent{value: 42}})
	history := es.Find(uuid)
	c.Assert(history.events, HasLen, 1)
	c.Assert(history.version, Equals, 1)
}
