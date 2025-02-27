package commands

import (
	"errors"
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"roman/port"
	"roman/util"
	"roman/util/enum"
	"strings"
	"time"
)

type NearBd struct {
	birthdayService port.BirthdayService
}

func (s NearBd) Handle(e any) util.RomanError {
	interaction, ok := e.(*events.ApplicationCommandInteractionCreate)
	if !ok {
		return util.NewErrorWithDisplay("[NearBd Handler]", errors.New("failed to convert discord event to ApplicationCommandInteractionCreate"), "Couldn't find embed")
	}

	currentTime := time.Now()

	var birthdayLimit uint = 15
	date := uint(currentTime.Month()*100) + uint(currentTime.Day())
	birthdays, err := s.birthdayService.GetBirthdaysFromDate(date, 1, birthdayLimit)
	if err != nil {
		message := discord.NewMessageCreateBuilder().
			SetContent(err.DisplayError()).
			Build()

		return util.NewErrorWithDisplay("[NearBd Handler]", interaction.CreateMessage(message), "Failed to send birthday fail ack")
	}

	var messageContent strings.Builder
	todayZeroed := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, time.Local)
	for _, birthday := range birthdays {
		date := time.Date(time.Now().Year(), time.Month(birthday.Date/100), int(birthday.Date%100), 12, 0, 0, 0, time.UTC)
		if todayZeroed.After(date) {
			date = date.AddDate(1, 0, 0)
		}

		messageContent.WriteString(fmt.Sprintf("<@%d>: <t:%d:D>", birthday.DiscordId, date.Unix()))
		if birthday.BirthYear != 0 {
			messageContent.WriteString(fmt.Sprintf(" (%d)", uint(date.Year())-birthday.BirthYear))
		}

		messageContent.WriteRune('\n')
	}

	totalBirthdayCount, err := s.birthdayService.GetTotalBirthdayCount()

	btnNext := discord.NewSecondaryButton("▶", fmt.Sprintf("%s-%s", enum.ChangeBdEmbedPaginator, &util.Paginator{
		CreatorId: uint64(interaction.User().ID),
		Action:    "next",
		Page:      1,
	}))
	if totalBirthdayCount == len(birthdays) {
		btnNext = btnNext.AsDisabled()
	}

	btnPrev := discord.NewSecondaryButton("◀", fmt.Sprintf("%s-%s", enum.ChangeBdEmbedPaginator, &util.Paginator{
		CreatorId: uint64(interaction.User().ID),
		Action:    "prev",
		Page:      1,
	})).AsDisabled()

	embed := discord.
		NewEmbedBuilder().
		SetTitle(":white_flower: Ближайшие дни рождения :white_flower:").
		SetDescription(messageContent.String()).
		SetColor(0xDB20E8).
		Build()

	message := discord.NewMessageCreateBuilder().
		SetAllowedMentions(&discord.AllowedMentions{
			Parse:       nil,
			Roles:       nil,
			Users:       nil,
			RepliedUser: false,
		}).
		SetEmbeds(embed).
		AddActionRow(btnPrev, btnNext).
		Build()
	return util.NewErrorWithDisplay("[NearBd Handler]", interaction.CreateMessage(message), "Failed to send birthday list")
}

func (s NearBd) Info() discord.SlashCommandCreate {
	return discord.SlashCommandCreate{
		Name:        SlashNearBd,
		Description: "Show nearest birthdays",
	}
}
