package service

import (
	"roman/env"
	"roman/util"
	"strconv"
)

type ConfigService struct {
	discordToken    string
	discordGuildId  uint64
	osuClientId     string
	osuClientSecret string
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

func NewConfigService() *ConfigService {
	guildId, _ := strconv.Atoi(util.GetEnv(env.GUILD_ID))

	return &ConfigService{
		discordToken:    util.GetEnv(env.BOT_TOKEN),
		discordGuildId:  uint64(guildId),
		osuClientId:     util.GetEnv(env.OSU_API_CLIENT_ID),
		osuClientSecret: util.GetEnv(env.OSU_API_KEY),
	}
}
