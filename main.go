package main

import (
	"context"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"roman/commands"
	"roman/env"
	romanEvents "roman/events"
	"roman/util"
	"strconv"
	"strings"
	"syscall"
)

func main() {
	// Load env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Should probably make this a config service or something
	token := util.GetEnv(env.BOT_TOKEN)
	guildID, _ := strconv.Atoi(util.GetEnv(env.GUILD_ID))

	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuilds,
				gateway.IntentGuildMessages,
				gateway.IntentGuildMessageReactions,
			)),
		bot.WithEventListenerFunc(onSlashCommand),
		bot.WithEventListenerFunc(onInteractionComponent),
	)
	if err != nil {
		log.Fatalf("Couldn't create bot client: %v\r\n", err)
	}
	defer client.Close(context.TODO())

	// Get all command description objects
	var discordCommands []discord.ApplicationCommandCreate
	for _, value := range commands.GetAll() {
		discordCommands = append(discordCommands, value.Info())
	}

	// Register commands
	if _, err = client.Rest().SetGuildCommands(client.ApplicationID(), snowflake.ID(guildID), discordCommands); err != nil {
		log.Fatal("error while registering commands", err)
	}

	if err = client.OpenGateway(context.TODO()); err != nil {
		log.Fatalf("Couldn't connect bot: %v\r\n", err)
	}

	log.Println("CTRL-C to exit.")
	// Listen for CTRL-C and other stuff
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
}

// onSlashCommand fires when a slash command is used.
func onSlashCommand(e *events.ApplicationCommandInteractionCreate) {
	data := e.SlashCommandInteractionData()
	log.Println("Command used:", data.CommandName())

	// Should probably introduce an api so that I don't have to interact with the inner components
	// 2025.01.04: Guess I did realize?
	registeredCommands := commands.GetAll()

	var err util.RomanError

	// Honestly, kinda sketch. Wonder if separating the commands in their own packages is a better solution overall.
	// 2024.12.20: TODO: Coming back to this, not sure why I didn't realize, but just like a .Get() method is the way to go
	command, ok := registeredCommands[data.CommandName()]
	if !ok {
		// This path shouldn't really happen, but in case there are discord shenanigans
		message := discord.NewMessageCreateBuilder().SetContent("Как я отправляю это сообщение?").Build()
		err = util.NewError("[onSlashCommand]", e.CreateMessage(message))
	} else {
		err = command.Handler(e)
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

// onInteractionComponent uhh not sure what this fires on exactly, but so far just button presses on embeds
func onInteractionComponent(e *events.ComponentInteractionCreate) {
	// Since this is for buttons, they all have their own "id" system.
	// The way it's handled currently an event is assigned i.e. usersRoles and an id at the end so that it's trackable
	registeredEvents := romanEvents.GetAll()
	currentEventName := strings.Split(e.Data.CustomID(), "-")[0]

	var err util.RomanError
	event, ok := registeredEvents[currentEventName]
	if !ok {
		log.Println("couldn't find event name:", currentEventName)
		message := discord.NewMessageCreateBuilder().SetContent("Couldn't find event name").Build()
		err = util.NewError("[onInteractionComponent]", e.CreateMessage(message))
	} else {
		err = event.Handler(e)
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
