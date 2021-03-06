package account

import (
	"fmt"
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

func (s *Service) act(c Command, done chan *UUID) {
	switch v := c.(type) {
	case OpenAccount:
		event := AccountOpened{InitialBalance: v.InitialBalance}
		for result := range s.store.Save([]Event{event}) {
			done <- result
		}
	case Credit:
		event := AccountCredited{v.Uuid, v.Amount}
		done <- s.actionOnAccount(v.Uuid, event)
	case Debit:
		event := AccountDebited{v.Uuid, v.Amount}
		done <- s.actionOnAccount(v.Uuid, event)
	case CreditOnTransfer:
		event := AccountCreditedOnTransfer{v.To, v.Transaction, v.Amount, v.From, v.To}
		done <- s.actionOnAccount(v.To, event)
	case DebitOnTransfer:
		event := AccountDebitedOnTransfer{v.From, v.Transaction, v.Amount, v.From, v.To}
		done <- s.actionOnAccount(v.From, event)
	}
	close(done)
}

func (s *Service) actionOnAccount(uuid *UUID, event Event) *UUID {
	account := <-s.store.Find(uuid)
	update := Update{uuid, []Event{event}, account.Version}
	err := <-s.store.Update(update)
	if err != nil {
		fmt.Println("Update on account failed", err)
	}
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

func (s *Service) Find(uuid *UUID) chan Account {
	done := make(chan Account)
	go func() {
		done <- Build(<-s.store.Find(uuid))
		close(done)
	}()
	return done
}
