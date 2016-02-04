package transfer

import (
	. "github.com/go-check/check"
	. "github.com/nu7hatch/gouuid"
	. "mcveat/cqrs/event"
	. "mcveat/cqrs/store"
)

type MyModelSuite struct{}

var _ = Suite(&MyModelSuite{})

func (s *MyModelSuite) TestBuilTransferWithNoEvents(c *C) {
	history := History{Events: []Event{}}
	transfer := Build(history)
	c.Assert(transfer, Equals, Transfer{})
}

func (s *MyModelSuite) TestBuildCreatedTransfer(c *C) {
	transaction, _ := NewV4()
	from, _ := NewV4()
	to, _ := NewV4()
	events := []Event{TransferCreated{transaction, from, to, 100}}
	history := History{Events: events}
	transfer := Build(history)
	c.Assert(transfer.Uuid, Equals, transaction)
	c.Assert(transfer.From, Equals, from)
	c.Assert(transfer.To, Equals, to)
	c.Assert(transfer.Amount, Equals, 100)
	c.Assert(transfer.Status, Equals, Created)
}

func (s *MyModelSuite) TestBuildDebitedTransfer(c *C) {
	transaction, _ := NewV4()
	from, _ := NewV4()
	to, _ := NewV4()
	events := []Event{TransferCreated{transaction, from, to, 100}, TransferDebited{transaction, from, to, 100}}
	history := History{Events: events}
	transfer := Build(history)
	c.Assert(transfer.Status, Equals, Debited)
}

func (s *MyModelSuite) TestBuildCreditedNotDebitedTransfer(c *C) {
	transaction, _ := NewV4()
	from, _ := NewV4()
	to, _ := NewV4()
	events := []Event{TransferCreated{transaction, from, to, 100}, TransferCredited{transaction, from, to, 100}}
	history := History{Events: events}
	transfer := Build(history)
	c.Assert(transfer.Status, Equals, Created)
}

func (s *MyModelSuite) TestBuildDebitedAndCreditedTransfer(c *C) {
	transaction, _ := NewV4()
	from, _ := NewV4()
	to, _ := NewV4()
	events := make([]Event, 3)
	events[0] = TransferCreated{transaction, from, to, 100}
	events[1] = TransferDebited{transaction, from, to, 100}
	events[2] = TransferCredited{transaction, from, to, 100}
	history := History{Events: events}
	transfer := Build(history)
	c.Assert(transfer.Status, Equals, Completed)
}

func (s *MyModelSuite) TestBuildCreditedAndDebitedTransfer(c *C) {
	transaction, _ := NewV4()
	from, _ := NewV4()
	to, _ := NewV4()
	events := make([]Event, 3)
	events[0] = TransferCreated{transaction, from, to, 100}
	events[1] = TransferCredited{transaction, from, to, 100}
	events[2] = TransferDebited{transaction, from, to, 100}
	history := History{Events: events}
	transfer := Build(history)
	c.Assert(transfer.Status, Equals, Completed)
}
