package events

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	events2 "github.com/disgoorg/disgo/events"
	"log/slog"
	"roman/adapters/commands"
	"roman/util"
	"strings"
)

type ComponentInteractionCreate struct {
	commands map[string]commands.Handler
}

func NewComponentInteractionCreateHandler(buttonCommands map[string]commands.Handler) *ComponentInteractionCreate {
	return &ComponentInteractionCreate{
		commands: buttonCommands,
	}
}

func (c *ComponentInteractionCreate) handle(e *events2.ComponentInteractionCreate) {
	// Since this is for buttons, they all have their own "id" system.
	// The way it's handled currently a command is assigned i.e. usersRoles and an id at the end so that it's trackable
	registeredCommands := c.AssociatedHandlers()
	currentEventName := strings.Split(e.Data.CustomID(), "-")[0]

	var err util.RomanError
	command, ok := registeredCommands[currentEventName]
	if !ok {
		//log.Println("couldn't find associated command with event name:", currentEventName)
		//message := discord.NewMessageCreateBuilder().SetContent("Couldn't find command associated with that action").Build()
		//err = util.NewError("[ComponentInteractionCreate handle]", e.CreateMessage(message))
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

func (c *ComponentInteractionCreate) Register(client bot.Client) {
	client.EventManager().AddEventListeners(bot.NewListenerFunc(c.handle))
}

func (c *ComponentInteractionCreate) AssociatedHandlers() map[string]commands.Handler {
	return c.commands
}
