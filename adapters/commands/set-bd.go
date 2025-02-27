package commands

import (
	"errors"
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/goodsign/monday"
	"roman/port"
	"roman/util"
	"strconv"
	"strings"
	"time"
)

type SetBd struct {
	birthdayService port.BirthdayService
}

func (s SetBd) Info() discord.SlashCommandCreate {
	return discord.SlashCommandCreate{
		Name:        SlashSetBd,
		Description: "Set your birthday",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionString{
				Name:        "date",
				Description: "Дата др (dd.mm)",
				Required:    true,
			},
			discord.ApplicationCommandOptionInt{
				Name:        "year",
				Description: "Год рождение (yyyy)",
			},
		},
	}
}

func (s SetBd) Handle(e any) util.RomanError {
	interaction, ok := e.(*events.ApplicationCommandInteractionCreate)
	if !ok {
		return util.NewErrorWithDisplay("[SetBd Handler]", errors.New("failed to convert discord event to ApplicationCommandInteractionCreate"), "Couldn't find embed")
	}

	data := interaction.SlashCommandInteractionData()
	dateString := data.String("date")
	year := data.Int("year")

	if year < 0 || year > 9999 {
		message := discord.NewMessageCreateBuilder().
			SetContent("неее, давай годик нормальный :nerd:").
			SetEphemeral(true).
			Build()

		return util.NewErrorWithDisplay("[SetBd Handler]", interaction.CreateMessage(message), "Failed to send year parse error")
	}

	dateArr := strings.Split(dateString, ".")
	if len(dateArr) != 2 || len(dateArr[0]) != 2 || len(dateArr[1]) != 2 {
		message := discord.NewMessageCreateBuilder().
			SetContent("неее, давай дату нормальную :nerd:").
			SetEphemeral(true).
			Build()

		return util.NewErrorWithDisplay("[SetBd Handler]", interaction.CreateMessage(message), "Failed to send date creation error")
	}

	_, err := time.Parse(time.DateOnly, fmt.Sprintf("%04d-%s-%s", year, dateArr[1], dateArr[0]))
	if err != nil {
		message := discord.NewMessageCreateBuilder().
			SetContent("Перепроверь дату указанную").
			SetEphemeral(true).
			Build()

		return util.NewErrorWithDisplay("[SetBd Handler]", interaction.CreateMessage(message), "Failed to send date check error")
	}

	day, _ := strconv.Atoi(dateArr[0])
	month, _ := strconv.Atoi(dateArr[1])

	customErr := s.birthdayService.SetBd(uint(interaction.User().ID), uint(month*100+day), uint(year))
	if customErr != nil {
		message := discord.NewMessageCreateBuilder().
			SetContent(customErr.DisplayError()).
			SetEphemeral(true).
			Build()

		return util.NewError("[SetBd Handler]", interaction.CreateMessage(message))
	}

	dateTime := time.Date(0, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	message := discord.NewMessageCreateBuilder().
		SetContent(fmt.Sprintf("Установил твой день рождения на %s", monday.Format(dateTime, "02 January", monday.LocaleRuRU))).
		Build()

	return util.NewError("[SetBd Handler]", interaction.CreateMessage(message))
}
