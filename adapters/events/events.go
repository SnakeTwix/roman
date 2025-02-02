package events

import (
	"roman/util"
	"roman/util/enum"
)

var events = map[string]Event{
	enum.SelectMenuUserRolesId: TourneyRoles{},
}

func GetAll() map[string]Event {
	// uhh, add init logic if it ever needs to happen. Overall, this is fine currently
	return events
}

//func Get(name string) (Event, error) {
//	event, ok := events[name]
//	if !ok {
//
//	}
//}

type Event interface {
	// Handler generalizing is kind of impossible here. Handlers should check what event they get themselves
	Handler(event any) util.RomanError
}
