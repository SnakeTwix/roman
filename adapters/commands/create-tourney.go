package commands

import (
	"errors"
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/rest"
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
	}
}

func (c CreateTourney) Handler(e *events.ApplicationCommandInteractionCreate) util.RomanError {
	data := e.SlashCommandInteractionData()
	tourneyName := data.String("tourney_name")
	teamName := data.String("team_name")

	guilds := rest.NewGuilds(e.Client().Rest())
	guildId := *e.GuildID()

	guildChannels, baseErr := guilds.GetGuildChannels(guildId)
	if baseErr != nil {
		return util.NewErrorWithDisplay("[CreateTourney Handler]", baseErr, "Failed to fetch guild channels")
	}

	var err util.RomanError
	// Set the context at the end of the function when stuff exits
	defer func() {
		if err != nil {
			err.WriteCurrentContext("[CreateTourney Handler]")
		}
	}()

	// Role creation
	role, err := c.createRole(guilds, guildId, tourneyName, teamName)
	if err != nil {
		err.AddErrorCase(errors.New("role creation"))
		return err
	}

	// Category creation
	category, err := c.createCategory(guilds, guildId, guildChannels, role, tourneyName)
	if err != nil {
		err.AddErrorCase(errors.New("category creation"))
		return err
	}

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

	return util.NewErrorWithDisplay("[CreateTourney Handler]", e.CreateMessage(message), "Failed to send select menu for players")
}

func (c CreateTourney) createRole(guilds rest.Guilds, guildId snowflake.ID, tourneyName string, teamName string) (*discord.Role, util.RomanError) {
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
		return nil, util.NewErrorWithDisplay("[createRole]", err, fmt.Sprintf("Couldn't create role: `%s`", roleName))
	}

	return role, nil

}

func (c CreateTourney) createCategory(guilds rest.Guilds, guildId snowflake.ID, guildChannels []discord.GuildChannel, categoryRole *discord.Role, tourneyName string) (discord.GuildChannel, util.RomanError) {
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
		return nil, util.NewErrorWithDisplay("[createCategory]", err, fmt.Sprintf("Couldn't create category: `%s`", categoryName))
	}

	return category, nil
}

func (c CreateTourney) createTextChannel(guilds rest.Guilds, guildId snowflake.ID, category discord.GuildChannel, tourneyName string, subName string, position int) (*discord.GuildChannel, util.RomanError) {
	var channelName = fmt.Sprintf("%s-%s", tourneyName, subName)

	log.Println("Creating channel:", channelName)
	guildChannelCreate := discord.GuildTextChannelCreate{
		Name:     channelName,
		Position: position,
		ParentID: category.ID(),
	}

	channel, err := guilds.CreateGuildChannel(guildId, guildChannelCreate)
	if err != nil {
		return nil, util.NewErrorWithDisplay("[createTextChannel]", err, fmt.Sprintf("Couldn't create text channel: `%s`", channelName))
	}

	return &channel, nil
}

func (c CreateTourney) createVoiceChannel(guilds rest.Guilds, guildId snowflake.ID, category discord.GuildChannel, tourneyName string, position int) (*discord.GuildChannel, util.RomanError) {
	var channelName = fmt.Sprintf("%s-voice", tourneyName)

	log.Println("Creating voice:", channelName)
	guildChannelCreate := discord.GuildVoiceChannelCreate{
		Name:     channelName,
		Position: position,
		ParentID: category.ID(),
	}

	channel, err := guilds.CreateGuildChannel(guildId, guildChannelCreate)
	if err != nil {
		return nil, util.NewErrorWithDisplay("[createVoiceChannel]", err, fmt.Sprintf("Couldn't create voice: `%s`", channelName))
	}

	return &channel, nil
}
