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

	firstAccount := <-as.Act(account.OpenAccount{InitialBalance: 100})
	secondAccount := <-as.Act(account.OpenAccount{InitialBalance: 0})

	as.Act(account.Credit{Uuid: firstAccount, Amount: 300})
	as.Act(account.Debit{Uuid: firstAccount, Amount: 50})
	transfer := ts.Act(transfer.CreateTransfer{firstAccount, secondAccount, 125})

	time.Sleep(1000 * time.Millisecond)

	fmt.Println("Accounts:")
	fmt.Println("---------")
	fmt.Println(<-as.Find(firstAccount))
	fmt.Println(<-as.Find(secondAccount))
	fmt.Println("Transfers:")
	fmt.Println("----------")
	fmt.Println(<-ts.Find(<-transfer))
	fmt.Println("Events:")
	fmt.Println("-------")
	page := <-es.Events(0, 100)
	for _, event := range page.Events {
		fmt.Println(event)
	}
}
