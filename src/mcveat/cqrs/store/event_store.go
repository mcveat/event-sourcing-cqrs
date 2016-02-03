package store

import (
	. "github.com/nu7hatch/gouuid"
)

type EventStore map[UUID][]Event

func Empty() EventStore {
	return make(map[UUID][]Event)
}

func (es *EventStore) Save(events []Event) *UUID {
	uuid, err := NewV4()
	if err != nil {
		panic(err)
	}
	for i, e := range events {
		events[i] = e.SetUUID(uuid)
	}
	(*es)[(*uuid)] = events
	return uuid
}
