package store

import (
	"fmt"
	. "github.com/go-check/check"
	. "github.com/nu7hatch/gouuid"
	. "mcveat/cqrs/event"
	"testing"
)

func Test(t *testing.T) {
	TestingT(t)
}

type MySuite struct{}

var _ = Suite(&MySuite{})

type GenericEvent struct {
	uuid  *UUID
	value int
}

func (e GenericEvent) SetUUID(uuid *UUID) Event {
	e.uuid = uuid
	return e
}

func (e GenericEvent) String() string {
	return fmt.Sprint("{GenericEvent: value=%g}", e.value)
}

func (s *MySuite) TestSaveEventsInStore(c *C) {
	es := Empty()
	uuid := <-es.Save([]Event{GenericEvent{value: 42}, GenericEvent{value: -1}})
	c.Assert(uuid, NotNil)
	c.Assert(es.store, HasLen, 1)
	events := es.store[(*uuid)]
	c.Assert(events, HasLen, 2)
	c.Assert(events[0], Equals, GenericEvent{uuid: uuid, value: 42})
	c.Assert(events[1], Equals, GenericEvent{uuid: uuid, value: -1})
}

func (s *MySuite) TestSaveEventsInLog(c *C) {
	es := Empty()
	firstUuid := <-es.Save([]Event{GenericEvent{value: 1}, GenericEvent{value: 2}})
	secondUuid := <-es.Save([]Event{GenericEvent{value: 3}, GenericEvent{value: 4}})
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
	history := <-es.Find(randomUUID)
	c.Assert(history.Events, HasLen, 0)
	c.Assert(history.Version, Equals, 0)
}

func (s *MySuite) TestFind(c *C) {
	es := Empty()
	uuid := <-es.Save([]Event{GenericEvent{value: 42}})
	history := <-es.Find(uuid)
	c.Assert(history.Events, HasLen, 1)
	c.Assert(history.Version, Equals, 1)
}

func (s *MySuite) TestUpdateNotExisting(c *C) {
	es := Empty()
	randomUUID, _ := NewV4()
	update := Update{randomUUID, []Event{}, 1}
	err := <-es.Update(update)
	c.Assert(err, ErrorMatches, "Called update on entity that does not exists.*")
}

func (s *MySuite) TestUpdateOptimisticLockFailed(c *C) {
	es := Empty()
	uuid := <-es.Save([]Event{GenericEvent{value: 42}, GenericEvent{value: 43}})
	update := Update{uuid, []Event{}, 1}
	err := <-es.Update(update)
	c.Assert(err, ErrorMatches, "Optimistic lock failed on update.*")
}

func (s *MySuite) TestUpdate(c *C) {
	es := Empty()
	uuid := <-es.Save([]Event{GenericEvent{value: 42}, GenericEvent{value: 43}})
	update := Update{uuid, []Event{GenericEvent{value: 44}}, 2}
	err := <-es.Update(update)
	c.Assert(err, IsNil)
	c.Assert(es.store[(*uuid)], HasLen, 3)
	c.Assert(es.log, HasLen, 3)
}

func (s *MySuite) TestEventsOnEmptySet(c *C) {
	es := Empty()
	page := es.Events(0, 20)
	c.Assert(page.Offset, Equals, 0)
	c.Assert(page.Events, HasLen, 0)
}

func (s *MySuite) TestEventsBadOffset(c *C) {
	es := storeWithEvents(50)
	page := es.Events(-5, 20)
	c.Assert(page.Offset, Equals, 20)
}

func (s *MySuite) TestEventsBadBatchSize(c *C) {
	es := storeWithEvents(50)
	page := es.Events(0, -20)
	c.Assert(page.Offset, Equals, 10)
}

func (s *MySuite) TestEvents(c *C) {
	es := storeWithEvents(50)
	page := es.Events(10, 10)
	c.Assert(page.Offset, Equals, 20)
	c.Assert(page.Events, HasLen, 10)
}

func (s *MySuite) TestEventsOnBoundary(c *C) {
	es := storeWithEvents(15)
	page := es.Events(10, 10)
	c.Assert(page.Offset, Equals, 15)
	c.Assert(page.Events, HasLen, 5)
}

func storeWithEvents(n int) EventStore {
	es := Empty()
	for i := 0; i < n; i++ {
		<-es.Save([]Event{GenericEvent{value: i}})
	}
	return es
}
