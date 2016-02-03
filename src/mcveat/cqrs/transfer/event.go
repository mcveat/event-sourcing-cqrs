package transfer

import (
	"fmt"
	. "github.com/nu7hatch/gouuid"
	. "mcveat/cqrs/store"
)

type TransferCreated struct {
	uuid   *UUID
	from   *UUID
	to     *UUID
	amount int
}

func (e TransferCreated) String() string {
	return fmt.Sprint("{TransferCreated: from=", e.from, " to=", e.to, " amount=", e.amount, "}")
}

func (e TransferCreated) SetUUID(uuid *UUID) Event {
	e.uuid = uuid
	return e
}

type TransferDebited struct {
	uuid   *UUID
	from   *UUID
	to     *UUID
	amount int
}

func (e TransferDebited) String() string {
	return fmt.Sprint("{TransferDebited: from=", e.from, " to=", e.to, " amount=", e.amount, "}")
}

func (e TransferDebited) SetUUID(uuid *UUID) Event {
	e.uuid = uuid
	return e
}

type TransferCompleted struct {
	uuid   *UUID
	from   *UUID
	to     *UUID
	amount int
}

func (e TransferCompleted) String() string {
	return fmt.Sprint("{TransferCompleted: from=", e.from, " to=", e.to, " amount=", e.amount, "}")
}

func (e TransferCompleted) SetUUID(uuid *UUID) Event {
	e.uuid = uuid
	return e
}
