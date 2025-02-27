package events

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	events2 "github.com/disgoorg/disgo/events"
	"log"
	"log/slog"
	"roman/adapters/commands"
	"roman/util"
)

type ApplicationCommandInteractionCreate struct {
	commands map[string]commands.Handler
}

func NewApplicationCommandInteractionCreateHandler(commands map[string]commands.Handler) *ApplicationCommandInteractionCreate {
	return &ApplicationCommandInteractionCreate{
		commands: commands,
	}
}

func (c *ApplicationCommandInteractionCreate) handle(e *events2.ApplicationCommandInteractionCreate) {
	data := e.SlashCommandInteractionData()
	log.Println("Command used:", data.CommandName())

	// Should probably introduce an api so that I don't have to interact with the inner components
	// 2025.01.04: Guess I did realize?
	registeredCommands := c.AssociatedHandlers()

	var err util.RomanError

	// Honestly, kinda sketch. Wonder if separating the commands in their own packages is a better solution overall.
	// 2024.12.20: TODO: Coming back to this, not sure why I didn't realize, but just like a .Get() method is the way to go
	command, ok := registeredCommands[data.CommandName()]
	if !ok {
		// This path shouldn't really happen, but in case there are discord shenanigans
		message := discord.NewMessageCreateBuilder().SetContent("Как я отправляю это сообщение?").Build()
		err = util.NewError("[onSlashCommand]", e.CreateMessage(message))
	} else {
		err = command.Handle(e)
	}

	if err != nil {
		message := discord.NewMessageCreateBuilder().
			SetContent(err.DisplayError()).
			Build()

		// Don't feel like handling error fallback messages.
		// if THIS fails, might as well just kms
		_ = e.CreateMessage(message)

		e.Client().Logger().Error(err.DisplayError(), slog.Any("err", err))
	}
}

func (c *ApplicationCommandInteractionCreate) Register(client bot.Client) {
	client.EventManager().AddEventListeners(bot.NewListenerFunc(c.handle))
}

func (c *ApplicationCommandInteractionCreate) AssociatedHandlers() map[string]commands.Handler {
	return c.commands
}
