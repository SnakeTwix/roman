package port

type ConfigService interface {
	DiscordToken() string
	DiscordGuildId() uint64
	OsuClientId() string
	OsuClientSecret() string
	SqliteDbFile() string
}
