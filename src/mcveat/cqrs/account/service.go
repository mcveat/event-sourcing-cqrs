package account

import (
	. "github.com/nu7hatch/gouuid"
	. "mcveat/cqrs/store"
)

type Service struct {
	Store *EventStore
}

func (s *Service) Act(cmd Command) *UUID {
	switch v := cmd.(type) {
	case OpenAccount:
		event := AccountOpened{initialBalance: v.InitialBalance}
		return s.Store.Save([]Event{event})
	case Credit:
		account := s.Store.Find(v.Uuid)
		event := AccountCredited{uuid: v.Uuid, amount: v.Amount}
		update := Update{v.Uuid, []Event{event}, account.Version}
		s.Store.Update(update)
		return v.Uuid
	}
	return nil
}
