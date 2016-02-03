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
