package commands

import (
	"fmt"
	"roman/port"
	"roman/util"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type NearBd struct {
	birthdayService port.BirthdayService
}

func (s NearBd) Handler(e *events.ApplicationCommandInteractionCreate) util.RomanError {
	currentTime := time.Now()

	date := uint(currentTime.Month()*100) + uint(currentTime.Day())
	birthdays, err := s.birthdayService.GetBirthdaysFromDate(date, 100)
	if err != nil {
		message := discord.NewMessageCreateBuilder().
			SetContent(err.DisplayError()).
			Build()

		return util.NewErrorWithDisplay("[NearBd Handler]", e.CreateMessage(message), "Failed to send birthday fail ack")
	}

	var messageContent strings.Builder

	for _, birthday := range birthdays {
		date := time.Date(time.Now().Year(), time.Month(birthday.Date/100), int(birthday.Date%100), 12, 0, 0, 0, time.UTC)
		if time.Now().AddDate(0, 0, 1).After(date) {
			date = date.AddDate(1, 0, 0)
		}

		messageContent.WriteString(fmt.Sprintf("<@%d>: <t:%d:D>", birthday.DiscordId, date.Unix()))
		if birthday.BirthYear != 0 {
			messageContent.WriteString(fmt.Sprintf(" (%d)", uint(date.Year())-birthday.BirthYear))
		}

		messageContent.WriteRune('\n')
	}

	embed := discord.
		NewEmbedBuilder().
		SetTitle(":white_flower: Ближайшие дни рождения :white_flower:").
		SetDescription(messageContent.String()).
		SetColor(0xDB20E8).
		Build()

	//btnPrev := discord.NewSecondaryButton("", "back-id").WithEmoji(discord.ComponentEmoji{
	//	Name: "◀️",
	//})
	//
	//btnNext := discord.NewSecondaryButton("", "forward-id").WithEmoji(discord.ComponentEmoji{
	//	Name: "▶️",
	//})

	message := discord.NewMessageCreateBuilder().
		SetAllowedMentions(&discord.AllowedMentions{
			Parse:       nil,
			Roles:       nil,
			Users:       nil,
			RepliedUser: false,
		}).
		SetEmbeds(embed).
		//AddActionRow(btnPrev, btnNext).
		Build()
	return util.NewErrorWithDisplay("[NearBd Handler]", e.CreateMessage(message), "Failed to send birthday list")
}

func (s NearBd) Info() discord.SlashCommandCreate {
	return discord.SlashCommandCreate{
		Name:        NameNearBd,
		Description: "Show nearest birthdays",
	}
}
