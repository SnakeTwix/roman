package commands

import (
	"errors"
	api "github.com/SnakeTwix/gosu-api"
	"github.com/SnakeTwix/gosu-api/structs"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"roman/util"
	"strconv"
	"strings"
)

type ParseLobby struct {
	osuApi *api.Client
}

func (c ParseLobby) Info() discord.SlashCommandCreate {
	return discord.SlashCommandCreate{
		Name:        SlashParseLobby,
		Description: "Parses the multi lobby and spits out semicolon separated values of scores per player",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionString{
				Name:        "lobby_link",
				Description: "The link of the lobby where the maps were played",
				Required:    true,
			},
		},
	}
}

func (c ParseLobby) Handle(e any) util.RomanError {
	interaction, ok := e.(*events.ApplicationCommandInteractionCreate)
	if !ok {
		return util.NewErrorWithDisplay("[ParseLobby Handler]", errors.New("failed to convert discord event to ApplicationCommandInteractionCreate"), "Couldn't find embed")
	}

	data := interaction.SlashCommandInteractionData()
	lobbyLink := data.String("lobby_link")
	linkSplit := strings.Split(lobbyLink, "/")
	lobbyIdString := linkSplit[len(linkSplit)-1]

	lobbyId, err := strconv.Atoi(lobbyIdString)
	if err != nil {
		return util.NewErrorWithDisplay("[ParseLobby Handler]", errors.New("couldn't convert lobby id to integer"), "Couldn't parse lobby id")
	}

	matchData, err := c.osuApi.GetFullMatch(lobbyId)
	if matchData == nil {
		return util.NewErrorWithDisplay("[ParseLobby Handler]", err, "Couldn't fetch Match data")
	}

	type scoreData struct {
		mapId    int
		playerId int
	}

	// For easier playerIds lookup
	players := make(map[int]string)
	for _, user := range matchData.Users {
		players[user.Id] = user.Username
	}

	playerScores := make(map[scoreData]int)
	playedBeatmaps := make([]int, 0)

	// Note the scores of all the users playing a map
	for _, event := range matchData.Events {
		if event.Detail.Type != structs.MatchEventOther || len(event.Game.Scores) == 0 {
			continue
		}

		// Store each score as (playerId, beatmap) -> score
		// In the case of a map repetition, the last score is recorded
		for _, score := range event.Game.Scores {
			key := scoreData{
				mapId:    event.Game.BeatmapId,
				playerId: score.UserID,
			}
			playerScores[key] = score.Score
		}

		playedBeatmaps = append(playedBeatmaps, event.Game.BeatmapId)
	}

	var builder strings.Builder
	builder.WriteString("```")
	playerIds := make([]int, 0, len(players))

	// Get the playerIds and store them in an array as we need a consistent order
	for key := range players {
		builder.WriteString(players[key])
		builder.WriteRune(';')
		playerIds = append(playerIds, key)
	}
	builder.WriteString("\n")

	// For each beatmap, format a row of scores in the order of players
	for i := 0; i < len(playedBeatmaps); i++ {
		currentBeatmap := playedBeatmaps[i]
		for _, playerId := range playerIds {
			key := scoreData{
				mapId:    currentBeatmap,
				playerId: playerId,
			}

			score := strconv.Itoa(playerScores[key])
			builder.WriteString(score)
			builder.WriteString(";")
		}

		builder.WriteString("\n")
	}

	builder.WriteString("```")

	message := discord.NewMessageCreateBuilder().
		SetContent(builder.String()).
		Build()

	return util.NewErrorWithDisplay("[ParseLobby Handler]", interaction.CreateMessage(message), "Failed to send lobby stats")
}
