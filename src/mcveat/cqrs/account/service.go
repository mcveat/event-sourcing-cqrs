package account

import (
	. "github.com/nu7hatch/gouuid"
	. "mcveat/cqrs/store"
	"mcveat/cqrs/transfer"
	"time"
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
		event := AccountOpened{initialBalance: v.InitialBalance}
		return s.store.Save([]Event{event})
	case Credit:
		event := AccountCredited{v.Uuid, v.Amount}
		return s.actionOnAccount(v.Uuid, event)
	case Debit:
		event := AccountDebited{v.Uuid, v.Amount}
		return s.actionOnAccount(v.Uuid, event)
	case CreditOnTransfer:
		event := AccountCreditedOnTransfer{v.To, v.Transaction, v.Amount, v.From}
		return s.actionOnAccount(v.To, event)
	case DebitOnTransfer:
		event := AccountDebitedOnTransfer{v.From, v.Transaction, v.Amount, v.To}
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
	go s.listen()
}

func (s *Service) listen() {
	offset := 0
	for {
		time.Sleep(100 * time.Millisecond)
		page := s.store.Events(offset, 5)
		go s.handleEvents(page.Events)
		offset = page.Offset
	}
}

func (s *Service) handleEvents(events []Event) {
	for _, event := range events {
		s.handleEvent(event)
	}
}

func (s *Service) handleEvent(e Event) {
	switch event := e.(type) {
	case transfer.TransferCreated:
		s.Act(DebitOnTransfer{event.Uuid, event.Amount, event.From, event.To})
	}
}
