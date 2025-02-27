package commands

import (
	"errors"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"roman/util"
)

type Ping struct{}

func (c Ping) Info() discord.SlashCommandCreate {
	return discord.SlashCommandCreate{
		Name:        SlashPing,
		Description: "pong",
	}
}

func (c Ping) Handle(e any) util.RomanError {
	interaction, ok := e.(*events.ApplicationCommandInteractionCreate)
	if !ok {
		return util.NewErrorWithDisplay("[Ping Handler]", errors.New("failed to convert discord event to ApplicationCommandInteractionCreate"), "Couldn't find embed")
	}

	message := discord.NewMessageCreateBuilder().
		SetContent("Pong!").
		Build()

	return util.NewErrorWithDisplay("[Ping Handler]", interaction.CreateMessage(message), "Couldn't ping :frowning: ")
}
