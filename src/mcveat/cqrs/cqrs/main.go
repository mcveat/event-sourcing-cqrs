package main

import (
	"fmt"
	"mcveat/cqrs/store"
)

func main() {
	es := store.Empty()
	es.Save(1)
	fmt.Println(es)
}
