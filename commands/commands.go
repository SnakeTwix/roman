package commands

import "github.com/disgoorg/disgo/discord"

func GetAll() []discord.ApplicationCommandCreate {
	commands := []discord.ApplicationCommandCreate{
		ping,
	}

	return commands
}

const (
	NamePing = "ping"
)
