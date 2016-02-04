package store

import (
	"fmt"
	. "github.com/nu7hatch/gouuid"
)

type Event interface {
	SetUUID(uuid *UUID) Event
	String() string
}

type AccountOpened struct {
	Uuid           *UUID
	InitialBalance int
}

func (e AccountOpened) String() string {
	return fmt.Sprint("{AccountOpened: initialBalance=", e.InitialBalance, "}")
}

func (e AccountOpened) SetUUID(uuid *UUID) Event {
	e.Uuid = uuid
	return e
}

type AccountCredited struct {
	Uuid   *UUID
	Amount int
}

func (e AccountCredited) String() string {
	return fmt.Sprint("{AccountCredited: amount=", e.Amount, "}")
}

func (e AccountCredited) SetUUID(uuid *UUID) Event {
	e.Uuid = uuid
	return e
}

type AccountDebited struct {
	Uuid   *UUID
	Amount int
}

func (e AccountDebited) String() string {
	return fmt.Sprint("{AccountDebited: amount=", e.Amount, "}")
}

func (e AccountDebited) SetUUID(uuid *UUID) Event {
	e.Uuid = uuid
	return e
}

type AccountCreditedOnTransfer struct {
	Uuid        *UUID
	Transaction *UUID
	Amount      int
	From        *UUID
	To          *UUID
}

func (e AccountCreditedOnTransfer) String() string {
	return fmt.Sprint("{AccountCreditedOnTransfer: amount=", e.Amount, ", from=", e.From, "}")
}

func (e AccountCreditedOnTransfer) SetUUID(uuid *UUID) Event {
	e.Uuid = uuid
	return e
}

type AccountDebitedOnTransfer struct {
	Uuid        *UUID
	Transaction *UUID
	Amount      int
	From        *UUID
	To          *UUID
}

func (e AccountDebitedOnTransfer) String() string {
	return fmt.Sprint("{AccountDebitedOnTransfer: amount=", e.Amount, " to=", e.To, "}")
}

func (e AccountDebitedOnTransfer) SetUUID(uuid *UUID) Event {
	e.Uuid = uuid
	return e
}

type TransferCreated struct {
	Uuid   *UUID
	From   *UUID
	To     *UUID
	Amount int
}

func (e TransferCreated) String() string {
	return fmt.Sprint("{TransferCreated: from=", e.From, " to=", e.To, " amount=", e.Amount, "}")
}

func (e TransferCreated) SetUUID(uuid *UUID) Event {
	e.Uuid = uuid
	return e
}

type TransferDebited struct {
	Uuid   *UUID
	From   *UUID
	To     *UUID
	Amount int
}

func (e TransferDebited) String() string {
	return fmt.Sprint("{TransferDebited: from=", e.From, " to=", e.To, " amount=", e.Amount, "}")
}

func (e TransferDebited) SetUUID(uuid *UUID) Event {
	e.Uuid = uuid
	return e
}

type TransferCompleted struct {
	Uuid   *UUID
	From   *UUID
	To     *UUID
	Amount int
}

func (e TransferCompleted) String() string {
	return fmt.Sprint("{TransferCompleted: from=", e.From, " to=", e.To, " amount=", e.Amount, "}")
}

func (e TransferCompleted) SetUUID(uuid *UUID) Event {
	e.Uuid = uuid
	return e
}
