package cmd

import (
	"context"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/gateway"
	"github.com/go-co-op/gocron/v2"
	"log"
	"os"
	"os/signal"
	"roman/adapters/api/osu"
	"roman/adapters/commands"
	"roman/adapters/jobs"
	repo2 "roman/adapters/repo"
	"roman/service"
	"syscall"
)

func Start() error {
	configService := service.NewConfigService()

	// db
	db := repo2.InitDB(configService)

	// repo
	birthdayRepo := repo2.NewBirthdayRepo(db)

	// services
	birthdayService := service.NewBirthdayService(birthdayRepo)

	osuApi := osu.GetClient(configService)

	scheduler, err := gocron.NewScheduler()
	defer scheduler.Shutdown()
	if err != nil {
		return err
	}
	scheduler.Start()

	client, err := disgo.New(configService.DiscordToken(),
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuilds,
				gateway.IntentGuildMessages,
				gateway.IntentGuildMessageReactions,
			)),
	)

	if err != nil {
		log.Fatalf("Couldn't create bot client: %v\r\n", err)
	}

	discordCommands := commands.Init(osuApi, birthdayService)
	discordJobs := jobs.Init(configService, scheduler, client, birthdayRepo)
	discordBot := InitBot(configService, client, discordCommands, discordJobs)

	defer discordBot.client.Close(context.TODO())
	discordBot.Start()

	log.Println("CTRL-C to exit.")
	// Listen for CTRL-C and other stuff
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s

	return nil
}
