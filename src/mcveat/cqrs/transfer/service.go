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
	}
	return nil
}
