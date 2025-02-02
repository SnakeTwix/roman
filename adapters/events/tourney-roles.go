package events

import (
	"errors"
	"fmt"
	"github.com/disgoorg/disgo/discord"
	discordEvents "github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/rest"
	"roman/util"
	"strings"
)

type TourneyRoles struct{}

func (c TourneyRoles) Handler(e any) util.RomanError {
	interaction, ok := e.(*discordEvents.ComponentInteractionCreate)
	if !ok {
		return util.NewErrorWithDisplay("[TourneyRoles Handler]", errors.New("failed to convert discord event to ComponentInteractionCreate"), "Couldn't find embed")
	}

	data := interaction.UserSelectMenuInteractionData()
	roleId, _ := util.ParseUserSelectId(data.CustomID())

	guildId := *interaction.GuildID()
	members := rest.NewMembers(interaction.Client().Rest())

	var message strings.Builder
	message.WriteString("Selected users: ")

	for _, user := range data.Users() {
		err := members.AddMemberRole(guildId, user.ID, roleId)

		// TODO: Don't abort on failing to add a role. Try adding roles to all the remaining users and notify when a user fails to get a role
		if err != nil {
			errMessage := fmt.Sprintf("Couldn't add role to user: %s. Stopping Process", user.Username)
			return util.NewErrorWithDisplay("[TourneyRoles Handler]", err, errMessage)
		}

		message.WriteString(fmt.Sprintf("<@%d>", user.ID))
		message.WriteString("")
	}

	message.WriteString(": You've been added to a team!")
	discordMessage := discord.NewMessageCreateBuilder().SetContent(message.String()).Build()

	return util.NewErrorWithDisplay("[TourneyRoles Handler]", interaction.CreateMessage(discordMessage), "Failed to send message")
}
