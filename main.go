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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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

	var discordCommands []discord.ApplicationCommandCreate
	for _, value := range commands.GetAll() {
		discordCommands = append(discordCommands, value.Info())
	}

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

func onSlashCommand(e *events.ApplicationCommandInteractionCreate) {
	data := e.SlashCommandInteractionData()
	log.Println("Command used:", data.CommandName())

	registeredCommands := commands.GetAll()

	var err error

	// Honestly, kinda sketch. Wonder if separating the commands in their own packages is a better solution overall.
	// Have no idea how this would work multithreaded, which would be the ideal implementation
	command, ok := registeredCommands[data.CommandName()]
	if !ok {
		message := discord.NewMessageCreateBuilder().SetContent("Как ты сюда добрался??").Build()
		err = e.CreateMessage(message)
	} else {
		err = command.Handler(e)
	}

	if err != nil {
		e.Client().Logger().Error("error on sending response", slog.Any("err", err))
	}
}

func onInteractionComponent(e *events.ComponentInteractionCreate) {
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
