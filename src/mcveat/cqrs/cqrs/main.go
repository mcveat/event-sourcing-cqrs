package main

import (
	"fmt"
	"mcveat/cqrs/store"
)

func main() {
	es := store.Empty()
	es.Save([]store.Event{store.GenericEvent{}})
	fmt.Println(es)
}
