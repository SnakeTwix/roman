package main

import (
	"camar/commands"
	"camar/env"
	"camar/util"
	"context"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
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
	)
	if err != nil {
		log.Fatalf("Couldn't create bot client: %v\r\n", err)
	}
	defer client.Close(context.TODO())

	if _, err = client.Rest().SetGuildCommands(client.ApplicationID(), snowflake.ID(guildID), commands.GetAll()); err != nil {
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

	var err error
	switch data.CommandName() {
	case commands.NamePing:
		err = commands.Ping(e)
	}

	if err != nil {
		e.Client().Logger().Error("error on sending response", slog.Any("err", err))
	}
}
