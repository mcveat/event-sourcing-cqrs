package transfer

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

func (s *MySuite) TestCreateTransfer(c *C) {
	es := store.Empty()
	as := Service{&es}
	from, _ := NewV4()
	to, _ := NewV4()
	uuid := as.Act(CreateTransfer{from, to, 100})
	c.Assert(es.Events(0, 10).Events, HasLen, 1)
	history := es.Find(uuid)
	c.Assert(history.Events, HasLen, 1)
	c.Assert(history.Events[0], Equals, TransferCreated{uuid, from, to, 100})
	c.Assert(history.Version, Equals, 1)
}

func (s *MySuite) TestDebitAndCompleteTransfer(c *C) {
	es := store.Empty()
	as := Service{&es}
	from, _ := NewV4()
	to, _ := NewV4()
	uuid := as.Act(CreateTransfer{from, to, 100})
	as.Act(Debite{uuid, from, to, 100})
	as.Act(Complete{uuid, from, to, 100})
	c.Assert(es.Events(0, 10).Events, HasLen, 3)
	history := es.Find(uuid)
	c.Assert(history.Events, HasLen, 3)
	c.Assert(history.Events[0], Equals, TransferCreated{uuid, from, to, 100})
	c.Assert(history.Events[1], Equals, TransferDebited{uuid, from, to, 100})
	c.Assert(history.Events[2], Equals, TransferCompleted{uuid, from, to, 100})
	c.Assert(history.Version, Equals, 3)
}
