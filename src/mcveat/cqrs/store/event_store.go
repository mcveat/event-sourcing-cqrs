package store

import (
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
	for i, e := range events {
		events[i] = e.SetUUID(uuid)
	}
	es.store[(*uuid)] = events
	es.log = append(es.log, events...)
	return uuid
}

func (es *EventStore) Find(uuid *UUID) History {
	events, ok := es.store[(*uuid)]
	if !ok {
		events = make([]Event, 0)
	}
	return History{events, len(events)}
}
