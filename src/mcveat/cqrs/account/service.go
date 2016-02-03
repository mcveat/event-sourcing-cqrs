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
	account := s.Store.Find(uuid)
	update := Update{uuid, []Event{event}, account.Version}
	s.Store.Update(update)
	return uuid
}
