package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"roman/util"
)

var commands = map[string]Command{
	NamePing:          Ping{},
	NameCreateTourney: CreateTourney{},
	NameParseLobby:    ParseLobby{},
}

func GetAll() map[string]Command {
	// uhh, add init logic if it ever needs to happen. Overall, this is fine currently
	return commands
}

type Command interface {
	Info() discord.SlashCommandCreate
	Handler(*events.ApplicationCommandInteractionCreate) util.RomanError
}

const (
	NamePing          = "ping"
	NameCreateTourney = "create-tourney"
	NameParseLobby    = "parse-lobby"
)
