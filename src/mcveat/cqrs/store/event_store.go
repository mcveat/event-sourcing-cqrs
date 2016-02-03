package store

import (
	. "github.com/nu7hatch/gouuid"
)

type EventStore map[UUID]interface{}

func Empty() EventStore {
	return make(map[UUID]interface{})
}

func (es *EventStore) Save(event interface{}) *UUID {
	uuid, err := NewV4()
	if err != nil {
		panic(err)
	}
	(*es)[(*uuid)] = event
	return uuid
}
