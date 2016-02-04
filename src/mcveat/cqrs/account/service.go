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

func (s *Service) Act(cmd Command) chan *UUID {
	done := make(chan *UUID)
	go s.act(cmd, done)
	return done
}

func (s *Service) act(c Command, done chan *UUID) *UUID {
	var result *UUID
	switch v := c.(type) {
	case OpenAccount:
		event := AccountOpened{InitialBalance: v.InitialBalance}
		result = s.store.Save([]Event{event})
	case Credit:
		event := AccountCredited{v.Uuid, v.Amount}
		result = s.actionOnAccount(v.Uuid, event)
	case Debit:
		event := AccountDebited{v.Uuid, v.Amount}
		result = s.actionOnAccount(v.Uuid, event)
	case CreditOnTransfer:
		event := AccountCreditedOnTransfer{v.To, v.Transaction, v.Amount, v.From, v.To}
		result = s.actionOnAccount(v.To, event)
	case DebitOnTransfer:
		event := AccountDebitedOnTransfer{v.From, v.Transaction, v.Amount, v.From, v.To}
		result = s.actionOnAccount(v.From, event)
	}
	done <- result
	return result
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

func (s *Service) handleEvent(e Event) chan *UUID {
	switch event := e.(type) {
	case TransferCreated:
		return s.Act(DebitOnTransfer{event.Uuid, event.Amount, event.From, event.To})
	case AccountDebitedOnTransfer:
		return s.Act(CreditOnTransfer{event.Transaction, event.Amount, event.From, event.To})
	}
	return nil
}

func (s *Service) Find(uuid *UUID) Account {
	return Build(s.store.Find(uuid))
}
