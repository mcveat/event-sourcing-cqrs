package transfer

import (
	. "github.com/nu7hatch/gouuid"
	. "mcveat/cqrs/store"
)

type Service struct {
	store *EventStore
}

func NewService(es *EventStore) Service {
	return Service{es}
}

func (s *Service) Act(cmd Command) *UUID {
	switch v := cmd.(type) {
	case CreateTransfer:
		event := TransferCreated{from: v.From, to: v.To, amount: v.Amount}
		return s.store.Save([]Event{event})
	case Debite:
		event := TransferDebited{v.Uuid, v.From, v.To, v.Amount}
		return s.actionOnTransfer(v.Uuid, event)
	case Complete:
		event := TransferCompleted{v.Uuid, v.From, v.To, v.Amount}
		return s.actionOnTransfer(v.Uuid, event)
	}
	return nil
}

func (s *Service) actionOnTransfer(uuid *UUID, event Event) *UUID {
	account := s.store.Find(uuid)
	update := Update{uuid, []Event{event}, account.Version}
	s.store.Update(update)
	return uuid
}
