package account

import (
	. "github.com/go-check/check"
	. "github.com/nu7hatch/gouuid"
	"mcveat/cqrs/store"
	"testing"
)

func Test(t *testing.T) {
	TestingT(t)
}

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestOpenAccount(c *C) {
	es := store.Empty()
	as := Service{&es}
	uuid := as.Act(OpenAccount{InitialBalance: 100})
	c.Assert(es.Events(0, 10).Events, HasLen, 1)
	history := es.Find(uuid)
	c.Assert(history.Events, HasLen, 1)
	c.Assert(history.Events[0], Equals, AccountOpened{uuid: uuid, initialBalance: 100})
	c.Assert(history.Version, Equals, 1)
}

func (s *MySuite) TestCreditAccount(c *C) {
	es := store.Empty()
	as := Service{&es}
	uuid := as.Act(OpenAccount{InitialBalance: 100})
	as.Act(Credit{uuid, 200})
	c.Assert(es.Events(0, 10).Events, HasLen, 2)
	history := es.Find(uuid)
	c.Assert(history.Events, HasLen, 2)
	c.Assert(history.Events[0], Equals, AccountOpened{uuid: uuid, initialBalance: 100})
	c.Assert(history.Events[1], Equals, AccountCredited{uuid: uuid, amount: 200})
	c.Assert(history.Version, Equals, 2)
}

func (s *MySuite) TestDebitAccount(c *C) {
	es := store.Empty()
	as := Service{&es}
	uuid := as.Act(OpenAccount{InitialBalance: 100})
	as.Act(Debit{uuid, 50})
	c.Assert(es.Events(0, 10).Events, HasLen, 2)
	history := es.Find(uuid)
	c.Assert(history.Events, HasLen, 2)
	c.Assert(history.Events[0], Equals, AccountOpened{uuid: uuid, initialBalance: 100})
	c.Assert(history.Events[1], Equals, AccountDebited{uuid: uuid, amount: 50})
	c.Assert(history.Version, Equals, 2)
}

func (s *MySuite) TestCreditAccountOnTransfer(c *C) {
	es := store.Empty()
	as := Service{&es}
	uuid := as.Act(OpenAccount{InitialBalance: 100})
	transaction, _ := NewV4()
	otherAccount, _ := NewV4()
	as.Act(CreditOnTransfer{transaction, 200, otherAccount, uuid})
	c.Assert(es.Events(0, 10).Events, HasLen, 2)
	history := es.Find(uuid)
	c.Assert(history.Events, HasLen, 2)
	c.Assert(history.Events[0], Equals, AccountOpened{uuid: uuid, initialBalance: 100})
	c.Assert(history.Events[1], Equals, AccountCreditedOnTransfer{uuid, transaction, 200, otherAccount})
	c.Assert(history.Version, Equals, 2)
}

func (s *MySuite) TestDebitAccountOnTransfer(c *C) {
	es := store.Empty()
	as := Service{&es}
	uuid := as.Act(OpenAccount{InitialBalance: 100})
	transaction, _ := NewV4()
	otherAccount, _ := NewV4()
	as.Act(DebitOnTransfer{transaction, 50, uuid, otherAccount})
	c.Assert(es.Events(0, 10).Events, HasLen, 2)
	history := es.Find(uuid)
	c.Assert(history.Events, HasLen, 2)
	c.Assert(history.Events[0], Equals, AccountOpened{uuid: uuid, initialBalance: 100})
	c.Assert(history.Events[1], Equals, AccountDebitedOnTransfer{uuid, transaction, 50, otherAccount})
	c.Assert(history.Version, Equals, 2)
}
