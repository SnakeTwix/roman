package commands

import (
	"camar/enum"
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/snowflake/v2"
	"log"
	"strings"
)

type TourneyRoles struct{}

func (c TourneyRoles) Handler(e *events.ComponentInteractionCreate) error {
	data := e.UserSelectMenuInteractionData()
	guildId := *e.GuildID()
	roleId, _ := strings.CutPrefix(data.CustomID(), enum.SelectMenuUserRolesId)
	members := rest.NewMembers(e.Client().Rest())

	var message strings.Builder
	message.WriteString("Selected users: ")

	for _, user := range data.Users() {
		err := members.AddMemberRole(guildId, user.ID, snowflake.MustParse(roleId))
		if err != nil {
			log.Printf("Couldn't add role to user: %s %v", user.Username, err)
			return err
		}

		message.WriteString(fmt.Sprintf("<@%d>", user.ID))
		message.WriteString(" ")
	}

	discordMessage := discord.NewMessageCreateBuilder().SetContent(message.String()).Build()

	return e.CreateMessage(discordMessage)
}
