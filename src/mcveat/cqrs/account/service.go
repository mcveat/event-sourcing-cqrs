package account

import (
	. "github.com/nu7hatch/gouuid"
	"mcveat/cqrs/store"
)

type Service struct {
	Store *store.EventStore
}

func (s *Service) Act(cmd Command) *UUID {
	switch v := cmd.(type) {
	case OpenAccount:
		event := AccountOpened{initialBalance: v.InitialBalance}
		return s.Store.Save([]store.Event{event})
	}
	return nil
}
