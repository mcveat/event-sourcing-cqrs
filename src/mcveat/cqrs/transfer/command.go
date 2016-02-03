package transfer

import . "github.com/nu7hatch/gouuid"

type Command interface{}

type CreateTransfer struct {
	From   *UUID
	To     *UUID
	Amount int
}
