package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var ping = discord.SlashCommandCreate{
	Name:        NamePing,
	Description: "pong",
}

func Ping(e *events.ApplicationCommandInteractionCreate) error {
	_ = e.SlashCommandInteractionData()

	message := discord.NewMessageCreateBuilder().
		SetContent("pong").
		Build()

	return e.CreateMessage(message)
}
