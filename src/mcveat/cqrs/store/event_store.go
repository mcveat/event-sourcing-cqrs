package store

import (
	"fmt"
	. "github.com/nu7hatch/gouuid"
)

type EventStore struct {
	store map[UUID][]Event
	log   []Event
}

func Empty() EventStore {
	return EventStore{make(map[UUID][]Event), make([]Event, 0)}
}

func (es *EventStore) Save(events []Event) *UUID {
	uuid, err := NewV4()
	if err != nil {
		panic(err)
	}
	eventsWithParent := addUUID(uuid, events)
	es.store[(*uuid)] = eventsWithParent
	es.log = append(es.log, eventsWithParent...)
	return uuid
}

func (es *EventStore) Find(uuid *UUID) History {
	events, ok := es.store[(*uuid)]
	if !ok {
		events = make([]Event, 0)
	}
	return History{events, len(events)}
}

func (es *EventStore) Update(update Update) error {
	stored, ok := es.store[(*update.uuid)]
	if !ok {
		return fmt.Errorf("Called update on entity that does not exists: %g", update)
	}
	if len(stored) != update.version {
		return fmt.Errorf("Optimistic lock failed on update: %g", update)
	}
	eventsWithParent := addUUID(update.uuid, update.events)
	es.store[(*update.uuid)] = append(stored, eventsWithParent...)
	es.log = append(es.log, eventsWithParent...)
	return nil
}

func (es *EventStore) Events(offset int, batchSize int) Page {
	noOfEvents := len(es.log)
	if noOfEvents == 0 {
		return Page{0, make([]Event, 0)}
	}
	if offset < 0 {
		offset = 0
	}
	if offset >= noOfEvents {
		return Page{0, make([]Event, 0)}
	}
	if batchSize <= 0 {
		batchSize = 10
	}
	max := offset + batchSize
	if max > noOfEvents {
		max = noOfEvents
	}
	result := es.log[offset:max]
	return Page{offset + batchSize, result}
}

func addUUID(uuid *UUID, events []Event) []Event {
	result := make([]Event, 0, len(events))
	for _, e := range events {
		result = append(result, e.SetUUID(uuid))
	}
	return result
}
