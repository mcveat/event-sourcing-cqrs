package store

import (
	. "github.com/nu7hatch/gouuid"
	. "mcveat/cqrs/event"
)

type History struct {
	Events  []Event
	Version int
}

type Update struct {
	Uuid    *UUID
	Events  []Event
	Version int
}

type Page struct {
	Offset int
	Events []Event
}
