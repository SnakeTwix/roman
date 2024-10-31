package osu

import (
	api "github.com/SnakeTwix/gosu-api"
	"log"
	"roman/env"
	"roman/util"
	"strconv"
)

var osuClient *api.Client

func GetClient() *api.Client {
	if osuClient != nil {
		return osuClient
	}

	clientId, err := strconv.Atoi(util.GetEnv(env.OSU_API_CLIENT_ID))
	if err != nil {
		log.Fatalf("Provided clientId is not a number: %v\n", util.GetEnv(env.OSU_API_CLIENT_ID))
	}

	apiKey := util.GetEnv(env.OSU_API_KEY)

	client := api.New(clientId, apiKey)
	osuClient = &client

	err = osuClient.GetToken()
	if err != nil {
		log.Fatalf("Couldn't get a token: %s\n", err)
	}

	return osuClient
}
