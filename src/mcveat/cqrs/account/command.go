package account

import . "github.com/nu7hatch/gouuid"

type Command interface{}

type OpenAccount struct {
	InitialBalance int
}

type Credit struct {
	Uuid   *UUID
	Amount int
}

type Debit struct {
	Uuid   *UUID
	Amount int
}

type CreditOnTransfer struct {
	Transaction *UUID
	Amount      int
	From        *UUID
	To          *UUID
}

type DebitOnTransfer struct {
	Transaction *UUID
	Amount      int
	From        *UUID
	To          *UUID
}
