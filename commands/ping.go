package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"roman/util"
)

type Ping struct{}

func (c Ping) Info() discord.SlashCommandCreate {
	return discord.SlashCommandCreate{
		Name:        NamePing,
		Description: "pong",
	}
}

func (c Ping) Handler(e *events.ApplicationCommandInteractionCreate) util.RomanError {
	message := discord.NewMessageCreateBuilder().
		SetContent("Pong!").
		Build()

	return util.NewErrorWithDisplay("[Ping Handler]", e.CreateMessage(message), "Couldn't ping :frowning: ")
}
