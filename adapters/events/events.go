package events

import (
	"github.com/disgoorg/disgo/bot"
	"roman/adapters/commands"
)

func Register(events []Event, client bot.Client) {
	for _, event := range events {
		event.Register(client)
	}
}

type Event interface {
	AssociatedHandlers() map[string]commands.Handler
	Register(client bot.Client)
}
