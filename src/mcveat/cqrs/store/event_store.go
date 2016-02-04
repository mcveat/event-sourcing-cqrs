package store

import (
	"fmt"
	. "github.com/nu7hatch/gouuid"
	. "mcveat/cqrs/event"
)

type EventStore struct {
	store map[UUID][]Event
	log   []Event
}

func Empty() EventStore {
	return EventStore{make(map[UUID][]Event), make([]Event, 0)}
}

func (es *EventStore) Save(events []Event) chan *UUID {
	done := make(chan *UUID)
	go es.save(events, done)
	return done
}

func (es *EventStore) save(events []Event, done chan *UUID) {
	uuid, err := NewV4()
	if err != nil {
		panic(err)
	}
	eventsWithParent := addUUID(uuid, events)
	es.store[(*uuid)] = eventsWithParent
	es.log = append(es.log, eventsWithParent...)
	done <- uuid
	close(done)
}

func (es *EventStore) Find(uuid *UUID) chan History {
	done := make(chan History)
	go es.find(uuid, done)
	return done
}

func (es *EventStore) find(uuid *UUID, done chan History) {
	events, ok := es.store[(*uuid)]
	if !ok {
		events = make([]Event, 0)
	}
	done <- History{events, len(events)}
	close(done)
}

func (es *EventStore) Update(update Update) chan error {
	done := make(chan error)
	go es.update(update, done)
	return done
}

func (es *EventStore) update(update Update, done chan error) {
	stored, ok := es.store[(*update.Uuid)]
	if !ok {
		done <- fmt.Errorf("Called update on entity that does not exists: %g", update)
		close(done)
		return
	}
	if len(stored) != update.Version {
		done <- fmt.Errorf("Optimistic lock failed on update: %g", update)
		close(done)
		return
	}
	eventsWithParent := addUUID(update.Uuid, update.Events)
	es.store[(*update.Uuid)] = append(stored, eventsWithParent...)
	es.log = append(es.log, eventsWithParent...)
	done <- nil
  close(done)
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
		return Page{offset, make([]Event, 0)}
	}
	if batchSize <= 0 {
		batchSize = 10
	}
	max := offset + batchSize
	if max > noOfEvents {
		max = noOfEvents
	}
	result := es.log[offset:max]
	return Page{max, result}
}

func (es EventStore) String() string {
	return fmt.Sprint("{EventStore store=", es.store, ", log=", es.log, "}")
}

func addUUID(uuid *UUID, events []Event) []Event {
	result := make([]Event, 0, len(events))
	for _, e := range events {
		result = append(result, e.SetUUID(uuid))
	}
	return result
}
