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
	registeredCommands := commands.GetAll()

	var err error

	// Honestly, kinda sketch. Wonder if separating the commands in their own packages is a better solution overall.
	// Have no idea how this would work multithreaded, which would be the ideal implementation
	// 2024.12.20: TODO: Coming back to this, not sure why I didn't realize, but just like a .Get() method is the way to go
	command, ok := registeredCommands[data.CommandName()]
	if !ok {
		// This path shouldn't really happen, but in case there are discord shenanigans
		message := discord.NewMessageCreateBuilder().SetContent("Как ты сюда добрался??").Build()
		err = e.CreateMessage(message)
	} else {
		err = command.Handler(e)
	}

	// TODO: Introduce a custom error component and send message based on the error gotten from command
	if err != nil {
		e.Client().Logger().Error("error on sending response", slog.Any("err", err))
	}
}

// onInteractionComponent uhh not sure what this fires on exactly, but technically just button pressed on embeds
func onInteractionComponent(e *events.ComponentInteractionCreate) {
	// Since this is for buttons, they all have their own "id" system.
	// The way it's handled currently I assign an event i.e. usersRoles and an id at the end so that it's trackable
	registeredEvents := romanEvents.GetAll()
	currentEventName := strings.Split(e.Data.CustomID(), "-")[0]

	var err error
	event, ok := registeredEvents[currentEventName]
	if !ok {
		log.Println("couldn't find event name:", currentEventName)
		message := discord.NewMessageCreateBuilder().SetContent("Couldn't find event name").Build()
		err = e.CreateMessage(message)
	} else {
		err = event.Handler(e)
	}

	if err != nil {
		e.Client().Logger().Error("error on sending response", slog.Any("err", err))
	}
}
