package commands

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
	"log"
	"roman/util"
)

type CreateTourney struct{}

func (c CreateTourney) Info() discord.SlashCommandCreate {
	return discord.SlashCommandCreate{
		Name:        NameCreateTourney,
		Description: "Creates basic template for a tourney",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionString{
				Name:        "tourney_name",
				Description: "Tournament name abbreviation",
				Required:    true,
			},
			discord.ApplicationCommandOptionString{
				Name:        "team_name",
				Description: "Name of the team",
			},
		},
		DefaultMemberPermissions: json.NewNullablePtr(discord.PermissionAdministrator),
	}
}

func (c CreateTourney) Handler(e *events.ApplicationCommandInteractionCreate) error {
	data := e.SlashCommandInteractionData()
	tourneyName := data.String("tourney_name")
	teamName := data.String("team_name")

	guilds := rest.NewGuilds(e.Client().Rest())
	guildId := *e.GuildID()

	guildChannels, err := guilds.GetGuildChannels(guildId)
	if err != nil {
		log.Println("Failed to fetch guild channels", err)
		return err
	}

	// Role creation
	role, err := c.createRole(guilds, guildId, tourneyName, teamName)

	// Category creation
	category, err := c.createCategory(guilds, guildId, guildChannels, role, tourneyName)

	// Channel creation
	_, err = c.createTextChannel(guilds, guildId, category, tourneyName, "links", 1)
	if err != nil {
		return err
	}

	_, err = c.createTextChannel(guilds, guildId, category, tourneyName, "general", 2)
	if err != nil {
		return err
	}

	_, err = c.createTextChannel(guilds, guildId, category, tourneyName, "bots", 3)
	if err != nil {
		return err
	}

	_, err = c.createVoiceChannel(guilds, guildId, category, tourneyName, 4)
	if err != nil {
		return err
	}

	selectMenu := discord.NewUserSelectMenu(util.CreateUserSelectId(role.ID, e.User().ID), "Select users").
		WithMinValues(2).
		WithMaxValues(12)

	message := discord.NewMessageCreateBuilder().
		AddActionRow(selectMenu).
		SetEphemeral(true).
		Build()

	return e.CreateMessage(message)
}

func (c CreateTourney) createRole(guilds rest.Guilds, guildId snowflake.ID, tourneyName string, teamName string) (*discord.Role, error) {
	var roleName = tourneyName
	if teamName != "" {
		roleName = fmt.Sprintf("%s | %s", tourneyName, teamName)
	}

	log.Println("Creating role:", roleName)
	createRole := discord.RoleCreate{
		Name: roleName,
	}

	role, err := guilds.CreateRole(guildId, createRole)

	if err != nil {
		log.Println("Couldn't create role:", roleName, err)
		return nil, err
	}

	return role, err

}

func (c CreateTourney) createCategory(guilds rest.Guilds, guildId snowflake.ID, guildChannels []discord.GuildChannel, categoryRole *discord.Role, tourneyName string) (discord.GuildChannel, error) {
	var categoryName = tourneyName

	log.Println("Creating category:", categoryName)
	lastCategoryPosition := 0
	for i := 0; i < len(guildChannels); i++ {
		if guildChannels[i].Type() == discord.ChannelTypeGuildCategory {
			lastCategoryPosition = guildChannels[i].Position()
		}
	}

	guildCategoryCreate := discord.GuildCategoryChannelCreate{
		Name:     categoryName,
		Position: lastCategoryPosition + 1,

		// Make category private
		PermissionOverwrites: []discord.PermissionOverwrite{
			discord.RolePermissionOverwrite{
				// Everyone
				RoleID: guildId,
				Deny:   discord.PermissionViewChannel | discord.PermissionConnect,
			},

			// Created role for the category
			discord.RolePermissionOverwrite{
				RoleID: categoryRole.ID,
				Allow:  discord.PermissionViewChannel | discord.PermissionConnect,
			},
		},
	}

	category, err := guilds.CreateGuildChannel(guildId, guildCategoryCreate)
	if err != nil {
		log.Println("Couldn't create category:", categoryName, err)
		return nil, err
	}

	return category, err
}

func (c CreateTourney) createTextChannel(guilds rest.Guilds, guildId snowflake.ID, category discord.GuildChannel, tourneyName string, subName string, position int) (*discord.GuildChannel, error) {
	var channelName = fmt.Sprintf("%s-%s", tourneyName, subName)

	log.Println("Creating channel:", channelName)
	guildChannelCreate := discord.GuildTextChannelCreate{
		Name:     channelName,
		Position: position,
		ParentID: category.ID(),
	}

	channel, err := guilds.CreateGuildChannel(guildId, guildChannelCreate)
	if err != nil {
		log.Println("Couldn't create channel:", channelName, err)
		return nil, err
	}

	return &channel, err
}

func (c CreateTourney) createVoiceChannel(guilds rest.Guilds, guildId snowflake.ID, category discord.GuildChannel, tourneyName string, position int) (*discord.GuildChannel, error) {
	var channelName = fmt.Sprintf("%s-voice", tourneyName)

	log.Println("Creating voice:", channelName)
	guildChannelCreate := discord.GuildVoiceChannelCreate{
		Name:     channelName,
		Position: position,
		ParentID: category.ID(),
	}

	channel, err := guilds.CreateGuildChannel(guildId, guildChannelCreate)
	if err != nil {
		log.Println("Couldn't create voice:", channelName, err)
		return nil, err
	}

	return &channel, err
}
