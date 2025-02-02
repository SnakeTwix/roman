package service

import (
	"roman/util"
	"roman/util/env"
	"strconv"
)

type ConfigService struct {
	discordToken                     string
	discordGuildId                   uint64
	osuClientId                      string
	osuClientSecret                  string
	sqliteDbFile                     string
	discordBirthdayGreetingChannelId int
}

func (c *ConfigService) DiscordToken() string {
	return c.discordToken
}

func (c *ConfigService) DiscordGuildId() uint64 {
	return c.discordGuildId
}

func (c *ConfigService) OsuClientId() string {
	return c.osuClientId
}

func (c *ConfigService) OsuClientSecret() string {
	return c.osuClientSecret
}

func (c *ConfigService) SqliteDbFile() string {
	return c.sqliteDbFile
}

func (c *ConfigService) DiscordBirthdayGreetingChannelId() int {
	return c.discordBirthdayGreetingChannelId
}

func NewConfigService() *ConfigService {
	guildId, _ := strconv.Atoi(util.GetEnv(env.GUILD_ID))
	discordBirthdayGreetingChannelId, _ := strconv.Atoi(util.GetEnv(env.BIRTHDAY_GREETING_CHANNEL_ID))

	return &ConfigService{
		discordToken:                     util.GetEnv(env.BOT_TOKEN),
		discordGuildId:                   uint64(guildId),
		osuClientId:                      util.GetEnv(env.OSU_API_CLIENT_ID),
		osuClientSecret:                  util.GetEnv(env.OSU_API_KEY),
		sqliteDbFile:                     util.GetEnv(env.SQ_LITE_DB_FILE),
		discordBirthdayGreetingChannelId: discordBirthdayGreetingChannelId,
	}
}
