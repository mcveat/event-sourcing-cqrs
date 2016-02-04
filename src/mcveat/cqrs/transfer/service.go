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

func (s *Service) Act(cmd Command) chan *UUID {
	done := make(chan *UUID)
	go s.act(cmd, done)
	return done
}

func (s *Service) act(cmd Command, done chan *UUID) {
	switch v := cmd.(type) {
	case CreateTransfer:
		event := TransferCreated{From: v.From, To: v.To, Amount: v.Amount}
		for result := range s.store.Save([]Event{event}) {
			done <- result
		}
	case Debite:
		event := TransferDebited{v.Uuid, v.From, v.To, v.Amount}
		done <- s.actionOnTransfer(v.Uuid, event)
	case Complete:
		event := TransferCredited{v.Uuid, v.From, v.To, v.Amount}
		done <- s.actionOnTransfer(v.Uuid, event)
	}
	close(done)
}

func (s *Service) actionOnTransfer(uuid *UUID, event Event) *UUID {
	account := <-s.store.Find(uuid)
	update := Update{uuid, []Event{event}, account.Version}
	s.store.Update(update)
	return uuid
}

func (s *Service) StartListener() {
	go listener.Listen(s.store, s.handleEvent)
}

func (s *Service) handleEvent(e Event) chan *UUID {
	switch event := e.(type) {
	case AccountDebitedOnTransfer:
		return s.Act(Debite{event.Transaction, event.From, event.To, event.Amount})
	case AccountCreditedOnTransfer:
		return s.Act(Complete{event.Transaction, event.From, event.To, event.Amount})
	}
	return nil
}

func (s *Service) Find(uuid *UUID) chan Transfer {
	done := make(chan Transfer)
	go func() {
		done <- Build(<-s.store.Find(uuid))
		close(done)
	}()
	return done
}
