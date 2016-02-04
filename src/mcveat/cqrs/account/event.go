package account

import (
	"fmt"
	. "github.com/nu7hatch/gouuid"
	. "mcveat/cqrs/store"
)

type AccountOpened struct {
	uuid           *UUID
	initialBalance int
}

func (e AccountOpened) String() string {
	return fmt.Sprint("{AccountOpened: initialBalance=", e.initialBalance, "}")
}

func (e AccountOpened) SetUUID(uuid *UUID) Event {
	e.uuid = uuid
	return e
}

type AccountCredited struct {
	uuid   *UUID
	amount int
}

func (e AccountCredited) String() string {
	return fmt.Sprint("{AccountCredited: amount=", e.amount, "}")
}

func (e AccountCredited) SetUUID(uuid *UUID) Event {
	e.uuid = uuid
	return e
}

type AccountDebited struct {
	uuid   *UUID
	amount int
}

func (e AccountDebited) String() string {
	return fmt.Sprint("{AccountDebited: amount=", e.amount, "}")
}

func (e AccountDebited) SetUUID(uuid *UUID) Event {
	e.uuid = uuid
	return e
}

type AccountCreditedOnTransfer struct {
	uuid        *UUID
	transaction *UUID
	amount      int
	from        *UUID
	to          *UUID
}

func (e AccountCreditedOnTransfer) String() string {
	return fmt.Sprint("{AccountCreditedOnTransfer: amount=", e.amount, ", from=", e.from, "}")
}

func (e AccountCreditedOnTransfer) SetUUID(uuid *UUID) Event {
	e.uuid = uuid
	return e
}

type AccountDebitedOnTransfer struct {
	uuid        *UUID
	transaction *UUID
	amount      int
	from        *UUID
	to          *UUID
}

func (e AccountDebitedOnTransfer) String() string {
	return fmt.Sprint("{AccountDebitedOnTransfer: amount=", e.amount, " to=", e.to, "}")
}

func (e AccountDebitedOnTransfer) SetUUID(uuid *UUID) Event {
	e.uuid = uuid
	return e
}
