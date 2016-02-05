package store

import (
	"fmt"
	. "github.com/nu7hatch/gouuid"
	. "mcveat/cqrs/event"
	"sync"
)

type EventStore struct {
	store map[UUID][]Event
	log   []Event
	mux   sync.Mutex
}

func Empty() EventStore {
	return EventStore{store: make(map[UUID][]Event), log: make([]Event, 0)}
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
	es.synchronous(func() {
		es.store[(*uuid)] = eventsWithParent
		es.log = append(es.log, eventsWithParent...)
	})
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
	go es.synchronous(func() {
		es.update(update, done)
	})
	return done
}

func (es *EventStore) update(update Update, done chan error) {
	defer close(done)
	stored, ok := es.store[(*update.Uuid)]
	if !ok {
		done <- fmt.Errorf("Called update on entity that does not exists:", update)
		return
	}
	if len(stored) != update.Version {
		done <- fmt.Errorf("Optimistic lock failed on update:", update)
		return
	}
	eventsWithParent := addUUID(update.Uuid, update.Events)
	es.store[(*update.Uuid)] = append(stored, eventsWithParent...)
	es.log = append(es.log, eventsWithParent...)
	done <- nil
}

func (es *EventStore) synchronous(f func()) {
	es.mux.Lock()
	f()
	es.mux.Unlock()
}

func (es *EventStore) Events(offset int, batchsize int) chan Page {
	done := make(chan Page)
	go es.events(offset, batchsize, done)
	return done
}

func (es *EventStore) events(offset int, batchSize int, done chan Page) {
	defer close(done)
	noOfEvents := len(es.log)
	if noOfEvents == 0 {
		done <- Page{0, make([]Event, 0)}
		return
	}
	if offset < 0 {
		offset = 0
	}
	if offset >= noOfEvents {
		done <- Page{offset, make([]Event, 0)}
		return
	}
	if batchSize <= 0 {
		batchSize = 10
	}
	max := offset + batchSize
	if max > noOfEvents {
		max = noOfEvents
	}
	result := es.log[offset:max]
	done <- Page{max, result}
}

func addUUID(uuid *UUID, events []Event) []Event {
	result := make([]Event, 0, len(events))
	for _, e := range events {
		result = append(result, e.SetUUID(uuid))
	}
	return result
}
