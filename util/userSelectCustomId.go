package util

import (
	"fmt"
	"github.com/disgoorg/snowflake/v2"
	"roman/enum"
	"strings"
)

// CreateUserSelectId Spits out the CustomId for tourney creation
func CreateUserSelectId(roleId snowflake.ID, userId snowflake.ID) string {
	return fmt.Sprintf("%s-%d-%d", enum.SelectMenuUserRolesId, roleId, userId)
}

// ParseUserSelectId Returns the roleId and the User for whom the interaction was intended
func ParseUserSelectId(customId string) (snowflake.ID, snowflake.ID) {
	parsed := strings.Split(customId, "-")
	roleId := snowflake.MustParse(parsed[1])
	userId := snowflake.MustParse(parsed[2])

	return roleId, userId
}
