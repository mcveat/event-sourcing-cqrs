package transfer

import (
	"fmt"
	. "github.com/nu7hatch/gouuid"
	. "mcveat/cqrs/event"
	. "mcveat/cqrs/store"
)

const (
	Created   = iota
	Debited   = iota
	Completed = iota
)

type Transfer struct {
	Uuid     *UUID
	From     *UUID
	To       *UUID
	Amount   int
	Status   int
	debited  bool
	credited bool
}

func (t Transfer) String() string {
	return fmt.Sprint("{Transfer: uuid=", t.Uuid, ", from=", t.From, " to=", t.To, " amount=", t.Amount,
		" status=", statusString(t.Status), "}")
}

func statusString(status int) string {
	switch status {
	case Created:
		return "created"
	case Debited:
		return "debited"
	case Completed:
		return "completed"
	}
	return "unknown"
}

func Build(history History) Transfer {
	transfer := Transfer{}
	for _, event := range history.Events {
		transfer = apply(event, transfer)
	}
	return transfer
}

func apply(e Event, transfer Transfer) Transfer {
	switch event := e.(type) {
	case TransferCreated:
		transfer.Uuid = event.Uuid
		transfer.From = event.From
		transfer.To = event.To
		transfer.Amount = event.Amount
		transfer.Status = Created
	case TransferDebited:
		transfer.debited = true
	case TransferCredited:
		transfer.credited = true
	}
	if transfer.debited && transfer.credited {
		transfer.Status = Completed
	} else if transfer.debited {
		transfer.Status = Debited
	}
	return transfer
}
