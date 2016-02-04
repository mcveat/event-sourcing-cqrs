package listener

import (
	. "github.com/nu7hatch/gouuid"
	. "mcveat/cqrs/event"
	. "mcveat/cqrs/store"
	"time"
)

func Listen(es *EventStore, handle func(Event) chan *UUID) {
	offset := 0
	for {
		time.Sleep(100 * time.Millisecond)
		page := <-es.Events(offset, 5)
		for _, event := range page.Events {
			handle(event)
		}
		offset = page.Offset
	}
}
