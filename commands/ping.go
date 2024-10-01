package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type Ping struct{}

func (c Ping) Info() discord.SlashCommandCreate {
	return discord.SlashCommandCreate{
		Name:        NamePing,
		Description: "pong",
	}
}

func (c Ping) Handler(e *events.ApplicationCommandInteractionCreate) error {
	message := discord.NewMessageCreateBuilder().
		SetContent("это тут или че").
		Build()

	return e.CreateMessage(message)
}
