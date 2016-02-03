package account

import (
	. "github.com/go-check/check"
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
