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
