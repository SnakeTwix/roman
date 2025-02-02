package jobs

import (
	"fmt"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/go-co-op/gocron/v2"
	"log"
	"os"
	"path/filepath"
	"roman/port"
	"roman/util"
	"time"
)

func NewJobHappyBirthdayGreeting(configService port.ConfigService, scheduler gocron.Scheduler, client bot.Client, birthdayService port.BirthdayService) *JobHappyBirthdayGreeting {
	return &JobHappyBirthdayGreeting{
		configService:   configService,
		scheduler:       scheduler,
		client:          client,
		birthDayService: birthdayService,
	}
}

type JobHappyBirthdayGreeting struct {
	configService   port.ConfigService
	scheduler       gocron.Scheduler
	client          bot.Client
	birthDayService port.BirthdayRepo
}

func (j *JobHappyBirthdayGreeting) Register() util.RomanError {
	_, err := j.scheduler.NewJob(
		gocron.DailyJob(1,
			gocron.NewAtTimes(
				gocron.NewAtTime(
					5, 0, 0,
				),
			),
		),
		gocron.NewTask(j.Execute),
	)

	if err != nil {
		return util.NewError("[JobHappyBirthdayGreeting] Register", err)
	}

	return nil
}

func (j *JobHappyBirthdayGreeting) Execute() util.RomanError {
	date := time.Now()
	numberDate := int(date.Month())*100 + date.Day()
	birthdays, err := j.birthDayService.GetBirthdaysOnDate(uint(numberDate))
	if err != nil {
		log.Println(err)
		return err
	}

	if len(birthdays) == 0 {
		return nil
	}

	fiveSecTicker := time.NewTicker(time.Second * 5)
	defer fiveSecTicker.Stop()
	discordChannelId := snowflake.ID(j.configService.DiscordBirthdayGreetingChannelId())

	//* Should not error, but if it does then we're fucked
	currentPath, _ := os.Getwd()
	greetingImage, fileErr := os.Open(filepath.Join(currentPath, "/static/birthday_greeting.png"))
	if fileErr != nil {
		log.Println("[JobHappyBirthdayGreeting Execute]", "couldn't load greetingImage", fileErr)
		return util.NewError("[JobHappyBirthdayGreeting Execute]", fileErr)
	}

	for _, birthday := range birthdays {
		message := discord.NewMessageCreateBuilder().
			SetContent(fmt.Sprintf("<@%d> хихи :baby:", birthday.DiscordId)).
			AddFile("s-dr.png", "великий поздрав", greetingImage).
			Build()

		_, err := j.client.Rest().CreateMessage(discordChannelId, message)
		if err != nil {
			log.Println("[JobHappyBirthdayGreeting Execute]", err)
		}

		<-fiveSecTicker.C
	}

	return nil
}
