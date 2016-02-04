package account

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

func (s *Service) Act(c Command) *UUID {
	switch v := c.(type) {
	case OpenAccount:
		event := AccountOpened{InitialBalance: v.InitialBalance}
		return s.store.Save([]Event{event})
	case Credit:
		event := AccountCredited{v.Uuid, v.Amount}
		return s.actionOnAccount(v.Uuid, event)
	case Debit:
		event := AccountDebited{v.Uuid, v.Amount}
		return s.actionOnAccount(v.Uuid, event)
	case CreditOnTransfer:
		event := AccountCreditedOnTransfer{v.To, v.Transaction, v.Amount, v.From, v.To}
		return s.actionOnAccount(v.To, event)
	case DebitOnTransfer:
		event := AccountDebitedOnTransfer{v.From, v.Transaction, v.Amount, v.From, v.To}
		return s.actionOnAccount(v.From, event)
	}
	return nil
}

func (s *Service) actionOnAccount(uuid *UUID, event Event) *UUID {
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
	case TransferCreated:
		s.Act(DebitOnTransfer{event.Uuid, event.Amount, event.From, event.To})
	case AccountDebitedOnTransfer:
		s.Act(CreditOnTransfer{event.Transaction, event.Amount, event.From, event.To})
	}
}
