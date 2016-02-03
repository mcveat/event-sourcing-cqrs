package main

import (
	"fmt"
	"mcveat/cqrs/account"
	"mcveat/cqrs/store"
)

func main() {
	es := store.Empty()
	as := account.Service{&es}
	firstAccount := as.Act(account.OpenAccount{InitialBalance: 100})
	as.Act(account.Credit{Uuid: firstAccount, Amount: 300})
	fmt.Println(es)
}
