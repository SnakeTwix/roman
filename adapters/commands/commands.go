package commands

import (
	api "github.com/SnakeTwix/gosu-api"
	"github.com/disgoorg/disgo/discord"
	"roman/port"
	"roman/util"
)

type SlashCommands struct {
	commands map[string]SlashCommand
}

func InitSlashCommands(osuApi *api.Client, birthdayService port.BirthdayService) *SlashCommands {
	manager := SlashCommands{}

	commands := map[string]SlashCommand{
		SlashPing:          Ping{},
		SlashCreateTourney: CreateTourney{},
		SlashParseLobby: ParseLobby{
			// TODO: Add ports for future mocking
			osuApi: osuApi,
		},
		SlashSetBd: SetBd{
			birthdayService: birthdayService,
		},
		SlashNearBd: NearBd{
			birthdayService: birthdayService,
		},
	}

	manager.commands = commands

	return &manager
}

func (c *SlashCommands) GetAll() map[string]SlashCommand {
	return c.commands
}

func (c *SlashCommands) GetHandlers() map[string]Handler {
	handlers := map[string]Handler{}
	for name, command := range c.commands {
		handlers[name] = command
	}

	return handlers
}

type Handler interface {
	Handle(e any) util.RomanError
}

type Infoer interface {
	Info() discord.SlashCommandCreate
}

type SlashCommand interface {
	Handler
	Infoer
}

const (
	SlashPing          = "ping"
	SlashCreateTourney = "create-tourney"
	SlashParseLobby    = "parse-lobby"
	SlashSetBd         = "set-bd"
	SlashNearBd        = "near-bd"
)
