package util

import (
	"strings"
)

func CreateEventName(opts ...string) string {
	var builder strings.Builder

	for i := 0; i < len(opts); i++ {
		builder.WriteString(opts[i])

		if i != len(opts) {
			builder.WriteByte('-')
		}
	}

	return builder.String()
}
