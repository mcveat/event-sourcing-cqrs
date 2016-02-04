package account

import (
	. "github.com/go-check/check"
	. "github.com/nu7hatch/gouuid"
	. "mcveat/cqrs/event"
	. "mcveat/cqrs/store"
)

type MyModelSuite struct{}

var _ = Suite(&MyModelSuite{})

func (s *MyModelSuite) TestBuildAccountWithNoEvents(c *C) {
	history := History{Events: []Event{}}
	account := Build(history)
	c.Assert(account.Uuid, IsNil)
	c.Assert(account.Balance, Equals, 0)
}

func (s *MyModelSuite) TestBuildOpenedAccount(c *C) {
	uuid, _ := NewV4()
	events := []Event{AccountOpened{uuid, 100}}
	history := History{Events: events}
	account := Build(history)
	c.Assert(account.Uuid, Equals, uuid)
	c.Assert(account.Balance, Equals, 100)
}

func (s *MyModelSuite) TestBuildCreditedAccount(c *C) {
	uuid, _ := NewV4()
	events := []Event{AccountOpened{uuid, 100}, AccountCredited{uuid, 50}}
	history := History{Events: events}
	account := Build(history)
	c.Assert(account.Uuid, Equals, uuid)
	c.Assert(account.Balance, Equals, 150)
}

func (s *MyModelSuite) TestBuildDebitedAccount(c *C) {
	uuid, _ := NewV4()
	events := []Event{AccountOpened{uuid, 100}, AccountDebited{uuid, 40}}
	history := History{Events: events}
	account := Build(history)
	c.Assert(account.Uuid, Equals, uuid)
	c.Assert(account.Balance, Equals, 60)
}

func (s *MyModelSuite) TestBuildAccountCreditedOnTransfer(c *C) {
	uuid, _ := NewV4()
	randomOtherUUID, _ := NewV4()
	events := []Event{AccountOpened{uuid, 100}, AccountCreditedOnTransfer{uuid, randomOtherUUID, 50, randomOtherUUID, uuid}}
	history := History{Events: events}
	account := Build(history)
	c.Assert(account.Uuid, Equals, uuid)
	c.Assert(account.Balance, Equals, 150)
}

func (s *MyModelSuite) TestBuildAccountDebitedOnTransfer(c *C) {
	uuid, _ := NewV4()
	randomOtherUUID, _ := NewV4()
	events := []Event{AccountOpened{uuid, 100}, AccountDebitedOnTransfer{uuid, randomOtherUUID, 20, randomOtherUUID, uuid}}
	history := History{Events: events}
	account := Build(history)
	c.Assert(account.Uuid, Equals, uuid)
	c.Assert(account.Balance, Equals, 80)
}
