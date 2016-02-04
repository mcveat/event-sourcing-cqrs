package transfer

import (
	. "github.com/nu7hatch/gouuid"
	. "mcveat/cqrs/event"
	"mcveat/cqrs/listener"
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
		event := TransferCreated{From: v.From, To: v.To, Amount: v.Amount}
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

func (s *Service) StartListener() {
	go listener.Listen(s.store, s.handleEvent)
}

func (s *Service) handleEvent(e Event) {
	switch event := e.(type) {
	case AccountDebitedOnTransfer:
		s.Act(Debite{event.Transaction, event.From, event.To, event.Amount})
	}
}
