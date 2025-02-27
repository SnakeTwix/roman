package osu

import (
	api "github.com/SnakeTwix/gosu-api"
	"log"
	"roman/port"
	"roman/util"
	"roman/util/env"
	"strconv"
)

var osuClient *api.Client

func GetClient(config port.ConfigService) *api.Client {
	if osuClient != nil {
		return osuClient
	}

	clientId, err := strconv.Atoi(config.OsuClientId())
	if err != nil {
		log.Fatalf("Provided clientId is not a number: %v\n", util.GetEnv(env.OSU_API_CLIENT_ID))
	}

	apiKey := config.OsuClientSecret()

	client, err := api.NewClient(clientId, apiKey)
	if err != nil {
		log.Fatalf("Couldn't get a token: %s\n", err)
	}

	osuClient = &client
	return osuClient
}
