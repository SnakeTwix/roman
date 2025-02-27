package cmd

import (
	"context"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
	"github.com/go-co-op/gocron/v2"
	"log"
	"os"
	"os/signal"
	"roman/adapters/api/osu"
	"roman/adapters/commands"
	"roman/adapters/events"
	"roman/adapters/jobs"
	repo2 "roman/adapters/repo"
	"roman/service"
	"roman/util/enum"
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
	defer client.Close(context.TODO())

	discordJobs := jobs.Init(configService, scheduler, client, birthdayRepo)

	// Register jobs to be performed every X amount of time
	log.Println("Registering all jobs")

	for _, job := range discordJobs.GetAll() {
		// TODO: handle errors
		err := job.Register()
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println("Finished registering all jobs")

	// Register commands
	discordSlashCommands := commands.InitSlashCommands(osuApi, birthdayService)

	// Get all command description objects
	var discordSlashCommandsInfo []discord.ApplicationCommandCreate
	for _, value := range discordSlashCommands.GetAll() {
		discordSlashCommandsInfo = append(discordSlashCommandsInfo, value.Info())
	}

	if _, err := client.Rest().SetGuildCommands(client.ApplicationID(), snowflake.ID(configService.DiscordGuildId()), discordSlashCommandsInfo); err != nil {
		log.Fatal("error while registering commands", err)
	}

	// Register command listeners
	discordEventListeners := []events.Event{
		events.NewComponentInteractionCreateHandler(map[string]commands.Handler{
			enum.SelectMenuUserRolesId:  commands.TourneyRoles{},
			enum.ChangeBdEmbedPaginator: commands.NewUpdateBdList(birthdayService),
		}),
		events.NewApplicationCommandInteractionCreateHandler(discordSlashCommands.GetHandlers()),
	}

	log.Println("Registering event listeners")
	for _, eventListener := range discordEventListeners {
		eventListener.Register(client)
	}
	log.Println("Finished registering event listeners")

	if err := client.OpenGateway(context.TODO()); err != nil {
		log.Fatalf("Couldn't connect bot: %v\r\n", err)
	}

	log.Println("CTRL-C to exit.")
	// Listen for CTRL-C and other stuff
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s

	return nil
}
