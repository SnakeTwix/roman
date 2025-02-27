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

type UpdateBdList struct {
	birthdayService port.BirthdayService
}

func NewUpdateBdList(birthdayService port.BirthdayService) *UpdateBdList {
	return &UpdateBdList{birthdayService: birthdayService}
}

func (s UpdateBdList) Handle(e any) util.RomanError {
	interaction, ok := e.(*events.ComponentInteractionCreate)
	if !ok {
		return util.NewErrorWithDisplay("[UpdateBdList Handler]", errors.New("failed to convert discord event to ApplicationCommandInteractionCreate"), "Couldn't find embed")
	}

	eventData := strings.SplitN(interaction.Data.CustomID(), "-", 2)[1]
	paginatorData, err := util.ParsePaginator(eventData)
	if err != nil {
		return util.NewError("[UpdateBdList Handler]", err)
	}

	if paginatorData.Action == "prev" {
		paginatorData.Page--
	} else {
		paginatorData.Page++
	}

	currentTime := time.Now()

	var birthdayLimit uint = 15
	date := uint(currentTime.Month()*100) + uint(currentTime.Day())
	birthdays, err := s.birthdayService.GetBirthdaysFromDate(date, uint(paginatorData.Page), birthdayLimit)
	if err != nil {
		message := discord.NewMessageCreateBuilder().
			SetContent(err.DisplayError()).
			Build()

		return util.NewErrorWithDisplay("[UpdateBdList Handler]", interaction.CreateMessage(message), "Failed to send birthday fail ack")
	}

	message := interaction.Message
	var messageContent strings.Builder
	tomorrow := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day()+2, 0, 0, 0, 0, time.Local)
	for _, birthday := range birthdays {
		date := time.Date(time.Now().Year(), time.Month(birthday.Date/100), int(birthday.Date%100), 12, 0, 0, 0, time.UTC)
		if tomorrow.After(date) {
			date = date.AddDate(1, 0, 0)
		}

		messageContent.WriteString(fmt.Sprintf("<@%d>: <t:%d:D>", birthday.DiscordId, date.Unix()))
		if birthday.BirthYear != 0 {
			messageContent.WriteString(fmt.Sprintf(" (%d)", uint(date.Year())-birthday.BirthYear))
		}

		messageContent.WriteRune('\n')
	}
	message.Embeds[0].Description = messageContent.String()

	totalBirthdayCount, err := s.birthdayService.GetTotalBirthdayCount()

	btnNext := discord.NewSecondaryButton("▶", fmt.Sprintf("%s-%s", enum.ChangeBdEmbedPaginator, &util.Paginator{
		CreatorId: uint64(interaction.User().ID),
		Action:    "next",
		Page:      paginatorData.Page,
	}))

	lastPage := uint(totalBirthdayCount) / birthdayLimit
	if uint(totalBirthdayCount)%birthdayLimit != 0 {
		lastPage++
	}

	if lastPage == uint(paginatorData.Page) {
		btnNext = btnNext.AsDisabled()
	}

	btnPrev := discord.NewSecondaryButton("◀", fmt.Sprintf("%s-%s", enum.ChangeBdEmbedPaginator, &util.Paginator{
		CreatorId: uint64(interaction.User().ID),
		Action:    "prev",
		Page:      paginatorData.Page,
	}))
	if paginatorData.Page == 1 {
		btnPrev = btnPrev.AsDisabled()
	}

	updatedMessage := discord.NewMessageUpdateBuilder().
		SetEmbeds(message.Embeds...).
		AddActionRow(btnPrev, btnNext).
		Build()

	return util.NewErrorWithDisplay("[UpdateBdList Handler]", interaction.UpdateMessage(updatedMessage), "Failed to update birthday list")
}

func (s UpdateBdList) Info() discord.SlashCommandCreate {
	return discord.SlashCommandCreate{
		Name:        SlashNearBd,
		Description: "Show nearest birthdays",
	}
}
