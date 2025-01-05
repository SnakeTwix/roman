package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"roman/api/osu"
	"roman/commands"
	"roman/repo"
	"roman/service"
	"syscall"
)

func Start() {
	configService := service.NewConfigService()

	db := repo.InitDB(configService)

	birthdayRepo := repo.NewBirthdayRepo(db)

	birthdayService := service.NewBirthdayService(birthdayRepo)
	osuApi := osu.GetClient(configService)

	discordCommands := commands.Init(osuApi, birthdayService)
	discordBot := InitBot(configService, discordCommands)

	defer discordBot.client.Close(context.TODO())
	discordBot.Start()

	log.Println("CTRL-C to exit.")
	// Listen for CTRL-C and other stuff
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
}
