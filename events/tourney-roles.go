package events

import (
	"errors"
	"fmt"
	"github.com/disgoorg/disgo/discord"
	discordEvents "github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/rest"
	"log"
	"roman/util"
	"strings"
)

type TourneyRoles struct{}

func (c TourneyRoles) Handler(e any) error {
	interaction, ok := e.(*discordEvents.ComponentInteractionCreate)
	if !ok {
		return errors.New("wrong event data")
	}

	data := interaction.UserSelectMenuInteractionData()
	roleId, _ := util.ParseUserSelectId(data.CustomID())

	guildId := *interaction.GuildID()
	members := rest.NewMembers(interaction.Client().Rest())

	var message strings.Builder
	message.WriteString("Selected users: ")

	for _, user := range data.Users() {
		err := members.AddMemberRole(guildId, user.ID, roleId)
		if err != nil {
			log.Printf("Couldn't add role to user: %s %v", user.Username, err)
			return err
		}

		message.WriteString(fmt.Sprintf("<@%d>", user.ID))
		message.WriteString(" ")
	}

	discordMessage := discord.NewMessageCreateBuilder().SetContent(message.String()).Build()

	return interaction.CreateMessage(discordMessage)
}
