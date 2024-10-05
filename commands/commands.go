package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var commands = map[string]Command{
	NamePing:          Ping{},
	NameCreateTourney: CreateTourney{},
}

func GetAll() map[string]Command {
	// uhh, add init logic if it ever needs to happen. Overall, this is fine currently
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
