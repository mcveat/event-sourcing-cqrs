package account

import (
	"fmt"
	. "github.com/nu7hatch/gouuid"
	. "mcveat/cqrs/event"
	. "mcveat/cqrs/store"
)

type Account struct {
	Uuid    *UUID
	Balance int
}

func (a Account) String() string {
	return fmt.Sprint("{Account: uuid=", a.Uuid, " balance=", a.Balance, "}")
}

func Build(history History) Account {
	account := Account{}
	for _, event := range history.Events {
		account = apply(event, account)
	}
	return account
}

func apply(e Event, account Account) Account {
	switch event := e.(type) {
	case AccountOpened:
		account.Uuid = event.Uuid
		account.Balance = event.InitialBalance
	case AccountCredited:
		account.Balance += event.Amount
	case AccountCreditedOnTransfer:
		account.Balance += event.Amount
	case AccountDebited:
		account.Balance -= event.Amount
	case AccountDebitedOnTransfer:
		account.Balance -= event.Amount
	}
	return account
}
