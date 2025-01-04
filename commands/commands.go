package commands

import (
	api "github.com/SnakeTwix/gosu-api"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"roman/util"
)

type Commands struct {
	commands map[string]Command
}

func Init(osuApi *api.Client) *Commands {
	manager := Commands{}

	commands := map[string]Command{
		NamePing:          Ping{},
		NameCreateTourney: CreateTourney{},
		NameParseLobby: ParseLobby{
			// TODO: Add ports for future mocking
			osuApi: osuApi,
		},
	}

	manager.commands = commands

	return &manager
}

func (c *Commands) GetAll() map[string]Command {
	return c.commands
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
