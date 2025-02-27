package util

import (
	"errors"
	"fmt"
	"github.com/disgoorg/snowflake/v2"
	"strconv"
	"strings"
)

type Paginator struct {
	CreatorId uint64
	Action    string
	Page      uint64
}

func (p *Paginator) String() string {
	return fmt.Sprintf("%d-%d-%s", p.CreatorId, p.Page, p.Action)
}

func ParsePaginator(str string) (*Paginator, RomanError) {
	splitStr := strings.Split(str, "-")
	if len(splitStr) != 3 {
		return nil, NewErrorWithDisplay("[ParsePaginator]", errors.New("failed to split string"), "Could not parse pagination string")
	}

	paginator := Paginator{}

	creatorId, err := snowflake.Parse(splitStr[0])
	if err != nil {
		return nil, NewErrorWithDisplay("[ParsePaginator]", err, "failed to parse snowflake")
	}

	paginator.CreatorId = uint64(creatorId)

	paginator.Page, err = strconv.ParseUint(splitStr[1], 10, 64)
	if err != nil {
		return nil, NewErrorWithDisplay("[ParsePaginator]", err, "failed to parse Page number")
	}

	paginator.Action = splitStr[2]

	return &paginator, nil
}
