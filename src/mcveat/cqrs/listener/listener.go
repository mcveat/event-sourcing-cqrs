package listener

import (
	. "mcveat/cqrs/store"
	"time"
)

func Listen(es *EventStore, handle func(Event)) {
	offset := 0
	for {
		time.Sleep(100 * time.Millisecond)
		page := es.Events(offset, 5)
		for _, event := range page.Events {
			handle(event)
		}
		offset = page.Offset
	}
}
