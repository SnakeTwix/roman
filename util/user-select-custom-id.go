package util

import (
	"github.com/disgoorg/snowflake/v2"
	"roman/util/enum"
	"strings"
)

// CreateUserSelectId Spits out the CustomId for tourney creation player select
func CreateUserSelectId(roleId snowflake.ID, userId snowflake.ID) string {
	return CreateEventName(enum.SelectMenuUserRolesId, roleId.String(), userId.String())
}

// ParseUserSelectId Returns the roleId and the User for whom the interaction was intended
func ParseUserSelectId(customId string) (snowflake.ID, snowflake.ID) {
	parsed := strings.Split(customId, "-")

	// TODO: Decide if error handling is ever needed here
	roleId := snowflake.MustParse(parsed[1])
	userId := snowflake.MustParse(parsed[2])

	return roleId, userId
}
