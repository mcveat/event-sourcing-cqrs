package main

import (
	"fmt"
	"mcveat/cqrs/account"
	"mcveat/cqrs/store"
	"mcveat/cqrs/transfer"
	"time"
)

func main() {
	es := store.Empty()
	as := account.NewService(&es)
	ts := transfer.NewService(&es)

	as.StartListener()
	ts.StartListener()

	firstAccount := as.Act(account.OpenAccount{InitialBalance: 100})
	as.Act(account.Credit{Uuid: firstAccount, Amount: 300})
	as.Act(account.Debit{Uuid: firstAccount, Amount: 50})
	secondAccount := as.Act(account.OpenAccount{InitialBalance: 0})
	ts.Act(transfer.CreateTransfer{firstAccount, secondAccount, 125})

	time.Sleep(1000 * time.Millisecond)

	fmt.Println(es)
}
