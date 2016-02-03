package store

import (
	. "github.com/nu7hatch/gouuid"
)

type Event interface {
	SetUUID(uuid *UUID) Event
}

type History struct {
	events  []Event
	version int
}

type Update struct {
	uuid    *UUID
	events  []Event
	version int
}

type Page struct {
	offset int
	events []Event
}

type GenericEvent struct {
	uuid  *UUID
	value int
}

func (e GenericEvent) SetUUID(uuid *UUID) Event {
	e.uuid = uuid
	return e
}
