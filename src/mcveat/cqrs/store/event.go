package store

import . "github.com/nu7hatch/gouuid"

type Event interface {
	SetUUID(uuid *UUID) Event
	String() string
}

type History struct {
	Events  []Event
	Version int
}

type Update struct {
	uuid    *UUID
	events  []Event
	version int
}

type Page struct {
	Offset int
	Events []Event
}
