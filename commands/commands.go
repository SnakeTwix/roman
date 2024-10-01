package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func GetAll() []discord.ApplicationCommandCreate {
	commands := []discord.ApplicationCommandCreate{
		Ping{}.Info(),
		CreateTourney{}.Info(),
	}

	return commands
}

type Command interface {
	Info() discord.SlashCommandCreate
	Handler(*events.ApplicationCommandInteractionCreate) error
}

const (
	NamePing          = "ping"
	NameCreateTourney = "create-tourney"
)
